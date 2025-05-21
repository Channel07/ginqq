package ginqq

import (
	"github.com/gin-gonic/gin"
)

// RouterGroup 自定义路由组，强制处理函数接收自定义的*Context
type RouterGroup struct {
	*gin.RouterGroup
}

func (g *RouterGroup) GET(relativePath string, handlers ...func(*Context)) {
	g.handle("GET", relativePath, handlers)
}

func (g *RouterGroup) POST(relativePath string, handlers ...func(*Context)) {
	g.handle("POST", relativePath, handlers)
}
func (g *RouterGroup) DELETE(relativePath string, handlers ...func(*Context)) {
	g.handle("DELETE", relativePath, handlers)
}
func (g *RouterGroup) PUT(relativePath string, handlers ...func(*Context)) {
	g.handle("PUT", relativePath, handlers)
}

func (g *RouterGroup) handle(method, relativePath string, handlers []func(*Context)) {
	g.RouterGroup.Handle(method, relativePath, convertToGinHandlers(handlers)...)
}
