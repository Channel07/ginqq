# 个性化Gin框架（ginqq）
## 简介
本项目是一个基于Gin框架二次开发的框架，提供了一些部门通用的功能和扩展。
### 设计原则
- 开箱即用：内置系统级可观测性能力、API规范化能力，业务研发只需要关注业务逻辑
- 最小入侵：不更改或尽可能少更改暴露的公共方法，对业务代码极少入侵，历史项目无缝切换，无学习成本
- 功能可控：通过初始化参数控制各个功能开关，默认开启
- 可扩展性：不更改已有代码的基础上支持新增功能，如鉴权、限流

### 功能
| **功能类型**               | **功能清单**                                                                 | **说明**|                                                                 
|----------------------------|------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| **API规范化**              | - API规范化（服务端、客户端）<br>- 安全过滤| - API规范内容自动补全、不合规调用拦截<br>- 包含XSS、SQL注入、CSRF等安全过滤逻辑|
| **系统级可观测性**         | - 日志<br>- 监控<br>- 链路<br>        | - 内、外部流水<br>- 接口监控<br>- 接口链路  | 


## 快速开始
### 版本要求
- 推荐 Go 1.23+
### 安装
```bash
go get github.com/XXXX/ginqq

```
### 简单使用
默认配置
```go
package main

import (
	gin "chinatelecom.cn/framework/ginqq"
	"net/http"
)

func main() {
	r := gin.Default("Y122010101", "Y122")

	g := r.Group("/api")
	g.GET("/hello", func(c *gin.Context) {
		traceId := c.GetTraceId()
		http.Get("http://www.baidu.com")
		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!", "traceId": traceId})
	})

	r.Run(":8080")

}
```
个性化配置
```go
package main

import (
	gin "chinatelecom.cn/framework/ginqq"
	"net/http"
)

func main() {
	cfg := &gin.Config{
		ServiceCode: "Y122010101",
		PlatCode:    "Y122",
		MetricsConfig: &gin.MetricsConfig{
			Buckets: []float64{100, 200, 500, 1000, 2000, 3000, 5000, 10000},
			CustomLabels: map[string]string{
				"operate_type": "OPERATE_TYPE",
			},
		},
		HttpClientEnhanceConfig: &gin.HttpClientEnhanceConfig{
			Transport: http.DefaultTransport,
		},
	}
	r := gin.NewEngineWithConfig(cfg)

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
	})

	r.Run(":8080")

}
```

## 文档
- 部门API规范：[API规范文档](https://f9jctod099.feishu.cn/file/JuyabkP8RogqVyx4qo3cjnI6nEg)
- 部门日志规范：[日志规范文档](https://f9jctod099.feishu.cn/file/WhHzbmlSboIqI8xdwrJcADgYnVc)
- ginqq框架使用手册：[使用手册](https://f9jctod099.feishu.cn/docx/HqBtdOWukozvHlxnkvaccjhlnWd)



