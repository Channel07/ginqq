package ginqq

import (
	"errors"
	"fmt"
	"net/http"
)

var globalConfig *Config

type Config struct {
	// 基础信息，必传字段
	ServiceCode string // 服务编码
	PlatCode    string // 平台编码
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

// validate 配置校验逻辑，初始化默认值
func (c *Config) validate() error {
	var errs []error
	if globalConfig != nil {
		return errors.New("config already initialized")
	}
	// 必传字段校验
	if c.ServiceCode == "" {
		errs = append(errs, errors.New("ServiceCode is required, call WithServiceCode()"))
	}
	if c.PlatCode == "" {
		errs = append(errs, errors.New("PlatCode is required, call WithPlatCode()"))
	}

	if c.DisableMetrics == false {
		if c.MetricsConfig == nil { // 默认桶
			c.MetricsConfig = &MetricsConfig{
				Buckets: []float64{100, 200, 500, 1000, 2000, 3000, 5000, 10000}, // 默认桶
			}
		} else { // 自定义桶校验
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
		globalConfig = c
		return nil
	}
	return errors.Join(errs...)
}
