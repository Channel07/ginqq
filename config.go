package ginqq

import (
	"errors"
	"fmt"
	lumberjack "github.com/DeRuina/timberjack"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	cnf    *Config
	logger *logrus.Logger
)

type Config struct {
	SvcCode string // 服务编码（大写）
	AppName string // 应用名称（小写，以下划线拼接）

	// 功能开关&配置
	// 服务端系统可观测性
	DisableMetrics        bool
	MetricsConfig         *MetricsConfig
	DisableTracing        bool // 链路
	DisableTransactionLog bool // 内部流水

	// 服务端API规范化
	DisableApiStandardServer bool // 服务端API规范调用&校验拦截

	// Http客户端配置
	DisableHttpClientEnhance bool // http增强
	HttpClientEnhanceConfig  *HttpClientEnhanceConfig

	DisableProgramLog bool // 是否禁用程序日志

	LogConfig *LogConfig
}

type MetricsConfig struct {
	// bucket
	Buckets []float64 // 桶设置，单位毫秒，默认 [100, 200, 500, 1000, 2000, 3000, 5000, 10000]
	// 自定义标签,key:标签名，value:ctx中获取标签值的方法名，如：ctx.Get("operate_type")
	CustomLabels map[string]string
}

type LogConfig struct {
	// LogDir 是日志文件的目录，文件名自动生成。备份日志文件将保留在同一目录下。
	// 默认为 "/app/logs"（如果你的系统是 Windows 则默认为 "C:\\BllLogs\\<Config.SvcCode>_<Config.AppName>"）。
	LogDir string

	// MaxSize 是日志文件轮转前的最大大小（以兆字节为单位），默认为 1024MB。
	MaxSize int

	// MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	// 注意：一天定义为24小时，可能因夏令时、闰秒等因素与日历日不完全对应。
	// 默认不根据时间删除旧日志文件。
	MaxAge int

	// MaxBackups 是要保留的旧日志文件的最大数量，默认为 7。
	MaxBackups int

	// LocalTime 确定备份文件名中的时间戳是否使用计算机本地时间，默认使用 UTC 时间。
	LocalTime bool

	// Compress 确定轮转的日志文件是否使用 gzip 压缩，默认不压缩。
	Compress bool

	// RotationInterval 是日志轮转的最大时间间隔，默认为 1 天。
	// 如果自上次轮转经过的时间超过该间隔，即使文件大小未达到 MaxSize 也会触发轮转。
	// 最小推荐值为 1 分钟。如果设为 0 则禁用基于时间的轮转。
	//
	// 示例 RotationInterval = time.Hour * 24 表示每天轮转日志。
	RotationInterval time.Duration
}

type HttpClientEnhanceConfig struct {
	Transport                http.RoundTripper // 基础transport，默认使用http.DefaultTransport，可自定义transport设置连接池参数、超时时间等
	DisableSkipVerify        bool              // 跳过证书认证，默认跳过 TODO 后续增加证书认证体系
	DisableApiStandardClient bool              // 客户端API规范调用&校验拦截
	DisableTransactionLog    bool              // 外部流水
}

func (c *Config) GetPlayCode() string {
	return c.SvcCode[:4]
}

type PlainFormatter struct{}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

// LevelHook 分发不同级别的日志到不同文件
type LevelHook struct {
	writers   map[logrus.Level]io.Writer
	formatter logrus.Formatter
}

func (h *LevelHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0, len(h.writers))
	for level := range h.writers {
		levels = append(levels, level)
	}
	return levels
}

func (h *LevelHook) Fire(entry *logrus.Entry) error {
	writer, ok := h.writers[entry.Level]
	if !ok {
		return nil
	}
	msg, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	return err
}

// init 初始化默认配置。
func (c *Config) init() error {
	var errs []error
	if cnf != nil {
		return errors.New("config already initialized")
	}

	if c.SvcCode == "" {
		errs = append(errs, errors.New(`parameter "SvcCode" is required`))
	}
	if c.AppName == "" {
		errs = append(errs, errors.New(`parameter "AppName" is required`))
	}

	c.SvcCode = strings.ToUpper(strings.TrimSpace(c.SvcCode))
	c.AppName = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(c.AppName), "-", "_"))

	if !c.DisableMetrics {
		if c.MetricsConfig == nil {
			c.MetricsConfig = &MetricsConfig{
				Buckets: []float64{100, 200, 500, 1000, 2000, 3000, 5000, 10000}, // 默认桶
			}
		} else {
			if len(c.MetricsConfig.Buckets) == 0 {
				errs = append(errs, errors.New("MetricsConfig.Buckets cannot be empty"))
			} else {
				for i, b := range c.MetricsConfig.Buckets {
					if b <= 0 {
						errs = append(errs, fmt.Errorf("MetricsConfig.Buckets[%d] must > 0, got %f", i, b))
					}
				}
			}
		}
	}

	if c.LogConfig == nil {
		var logDir string
		if runtime.GOOS == "windows" {
			logDir = "C:\\BllLogs"
		} else {
			logDir = "/app/logs"
		}
		c.LogConfig = &LogConfig{
			LogDir:           logDir,
			MaxSize:          1024,
			MaxBackups:       7,
			RotationInterval: time.Hour * 24,
		}
	}

	if !c.DisableTransactionLog || !c.HttpClientEnhanceConfig.DisableTransactionLog {
		svcCode := strings.ToLower(c.SvcCode)
		filename := fmt.Sprintf(
			"%s/%s_%s/%s_%s_info-info.log",
			c.LogConfig.LogDir, svcCode, c.AppName, svcCode, c.AppName,
		)
		logger = &logrus.Logger{
			Out: &lumberjack.Logger{
				Filename:         filename,
				MaxSize:          c.LogConfig.MaxSize,
				MaxAge:           c.LogConfig.MaxAge,
				MaxBackups:       c.LogConfig.MaxBackups,
				LocalTime:        c.LogConfig.LocalTime,
				Compress:         c.LogConfig.Compress,
				RotationInterval: c.LogConfig.RotationInterval,
			},
			Formatter: new(PlainFormatter),
			Level:     logrus.InfoLevel,
		}
	}

	//if !c.DisableProgramLog {
	//	c.initProgramLog()
	//}

	if !c.DisableHttpClientEnhance && c.HttpClientEnhanceConfig == nil {
		c.HttpClientEnhanceConfig = &HttpClientEnhanceConfig{
			Transport:         http.DefaultTransport,
			DisableSkipVerify: false,
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	cnf = c
	return nil
}

func (c *Config) initProgramLog() {
	svcCode := strings.ToLower(c.SvcCode)

	logLevels := []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}

	writers := make(map[logrus.Level]io.Writer)

	traceFilename := fmt.Sprintf(
		"%s/%s_%s/trace/%s_%s_trace-trace.log",
		c.LogConfig.LogDir, svcCode, c.AppName, svcCode, c.AppName,
	)
	writers[logrus.TraceLevel] = &lumberjack.Logger{
		Filename:         traceFilename,
		MaxSize:          c.LogConfig.MaxSize,
		MaxAge:           c.LogConfig.MaxAge,
		MaxBackups:       c.LogConfig.MaxBackups,
		LocalTime:        c.LogConfig.LocalTime,
		Compress:         c.LogConfig.Compress,
		RotationInterval: c.LogConfig.RotationInterval,
	}
	debugFilename := fmt.Sprintf(
		"%s/%s_%s/debug/%s_%s_code-debug.log",
		c.LogConfig.LogDir, svcCode, c.AppName, svcCode, c.AppName,
	)
	writers[logrus.DebugLevel] = &lumberjack.Logger{
		Filename:         debugFilename,
		MaxSize:          c.LogConfig.MaxSize,
		MaxAge:           c.LogConfig.MaxAge,
		MaxBackups:       c.LogConfig.MaxBackups,
		LocalTime:        c.LogConfig.LocalTime,
		Compress:         c.LogConfig.Compress,
		RotationInterval: c.LogConfig.RotationInterval,
	}
	for _, level := range logLevels {
		filename := fmt.Sprintf(
			"%s/%s_%s/%s_%s_code-%s.log",
			c.LogConfig.LogDir, svcCode, c.AppName, svcCode, c.AppName, level,
		)
		writers[level] = &lumberjack.Logger{
			Filename:         filename,
			MaxSize:          c.LogConfig.MaxSize,
			MaxAge:           c.LogConfig.MaxAge,
			MaxBackups:       c.LogConfig.MaxBackups,
			LocalTime:        c.LogConfig.LocalTime,
			Compress:         c.LogConfig.Compress,
			RotationInterval: c.LogConfig.RotationInterval,
		}
	}

	logrus.AddHook(&LevelHook{writers, new(PlainFormatter)})
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(io.Discard)
}
