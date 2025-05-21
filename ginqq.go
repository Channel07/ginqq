package ginqq

import (
	"github.com/gin-gonic/gin"
)

type H gin.H

// GinQQ 自定义框架结构体。
type GinQQ struct {
	*gin.Engine
	Config *Config
}

// Group 返回自定义 RouterGroup 。
func (g *GinQQ) Group(relativePath string, handlers ...func(*Context)) *RouterGroup {
	return &RouterGroup{
		RouterGroup: g.Engine.Group(relativePath, convertToGinHandlers(handlers)...),
	}
}

func (g *GinQQ) GET(relativePath string, handlers ...func(*Context)) {
	g.Engine.GET(relativePath, convertToGinHandlers(handlers)...)
}

func (g *GinQQ) POST(relativePath string, handlers ...func(*Context)) {
	g.Engine.POST(relativePath, convertToGinHandlers(handlers)...)
}

func (g *GinQQ) PUT(relativePath string, handlers ...func(*Context)) {
	g.Engine.PUT(relativePath, convertToGinHandlers(handlers)...)
}

func (g *GinQQ) DELETE(relativePath string, handlers ...func(*Context)) {
	g.Engine.DELETE(relativePath, convertToGinHandlers(handlers)...)
}

func (g *GinQQ) Use(handlers ...func(*Context)) {
	g.Engine.Use(convertToGinHandlers(handlers)...)
}

func convertToGinHandlers(handlers []func(*Context)) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for i := range handlers {
		h := handlers[i]
		ginHandlers = append(ginHandlers, func(gc *gin.Context) {
			gc.Set(XMethodName, getRawHandlerName(h))
			h(Wrap(gc))
		})
	}
	return ginHandlers
}

// Default 创建默认的 GinQQ 实例。
func Default(svcCode, appName string) *GinQQ {
	return NewEngineWithConfig(&Config{
		SvcCode: svcCode,
		AppName: appName,
	})
}

func NewEngineWithConfig(config *Config) *GinQQ {
	if err := config.init(); err != nil {
		panic(err)
	}
	gq := &GinQQ{gin.New(), config}
	if !config.DisableTransactionLog {
		gq.Use(DispatchTransactionLog)
	}
	if !config.DisableHttpClientEnhance {
		HttpEnhance(config.HttpClientEnhanceConfig)
	}
	return gq
}
