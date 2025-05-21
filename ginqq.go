package ginqq

import (
	"github.com/gin-gonic/gin"
)

type H gin.H

// GinQQ 自定义框架结构体
type GinQQ struct {
	*gin.Engine
	Config *Config
}

// Group 重写分组，返回自定义RouterGroup
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

func (g *GinQQ) registerMiddleware() {
	//if g.Config.EnableInLog {
	//	g.Use(InJournalLogMiddleware())
	//}
	//if g.Config.EnableOutLog {
	//	g.Use(OutJournalLogMiddleware())
	//}

}

func convertToGinHandlers(handlers []func(*Context)) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, h := range handlers {
		ginHandlers = append(ginHandlers, func(gc *gin.Context) {
			h(Wrap(gc))
		})
	}
	return ginHandlers
}

// Default 创建默认的 GinQQ 实例
func Default(serviceCode string, platCode string) *GinQQ {
	cfg := &Config{
		ServiceCode: serviceCode,
		PlatCode:    platCode,
	}
	gq := NewEngineWithConfig(cfg)
	return gq
}

func (g *GinQQ) Use(handlers ...func(*Context)) {
	g.Engine.Use(convertToGinHandlers(handlers)...)
}

func NewEngineWithConfig(config *Config) *GinQQ {
	err := config.validate()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	gq := &GinQQ{
		Engine: r,
		Config: config,
	}
	// 注册中间件
	gq.registerMiddleware()
	// 启用HTTP增强
	if config.DisableHttpClientEnhance == false {
		HttpEnhance(config.HttpClientEnhanceConfig)
	}
	return gq
}
