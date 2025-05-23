package ginqq

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var cnf *Config

type Config struct {
	SvcCode string // 服务编码（大写）
	AppName string // 应用名称（小写，以下划线拼接）

	// 功能开关&配置
	// 服务端系统可观测性
	DisableMetrics       bool
	MetricsConfig        *MetricsConfig
	DisableTracing       bool // 链路
	DisableLoggingServer bool // 内部流水

	// 服务端API规范化
	DisableApiStandardServer bool // 服务端API规范调用&校验拦截

	// Http客户端配置
	DisableHttpClientEnhance bool // http增强
	HttpClientEnhanceConfig  *HttpClientEnhanceConfig
}

type MetricsConfig struct {
	// bucket
	Buckets []float64 // 桶设置，单位毫秒，默认 [100, 200, 500, 1000, 2000, 3000, 5000, 10000]
	// 自定义标签,key:标签名，value:ctx中获取标签值的方法名，如：ctx.Get("operate_type")
	CustomLabels map[string]string
}

type HttpClientEnhanceConfig struct {
	Transport                http.RoundTripper // 基础transport，默认使用http.DefaultTransport，可自定义transport设置连接池参数、超时时间等
	DisableSkipVerify        bool              // 跳过证书认证，默认跳过 TODO 后续增加证书认证体系
	DisableApiStandardClient bool              // 客户端API规范调用&校验拦截
	DisableLoggingClient     bool              // 外部流水
}

func (c *Config) GetPlayCode() string {
	return c.SvcCode[:4]
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

	c.SvcCode = strings.ToUpper(c.SvcCode)
	c.AppName = strings.ReplaceAll(strings.ToLower(c.AppName), "-", "_")

	if c.DisableMetrics == false {
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

	if c.DisableHttpClientEnhance == false && c.HttpClientEnhanceConfig == nil {
		c.HttpClientEnhanceConfig = &HttpClientEnhanceConfig{
			Transport:         http.DefaultTransport,
			DisableSkipVerify: false,
		}
	}

	if len(errs) == 0 {
		cnf = c
		return nil
	}
	return errors.Join(errs...)
}
