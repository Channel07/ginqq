// Package ginqq 自定义的context，用于传递系统上下文信息，如服务编码，traceId等
package ginqq

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

const (
	CtxKeyMethodCode = "Method-Code"
	CtxKeyTraceID    = "Trace-ID"
)

func Wrap(c *gin.Context) *Context {
	return &Context{Context: c}
}

// SetMethodCode 框架中间件数据传递
func (c *Context) SetMethodCode(methodCode string) {
	c.Set(CtxKeyMethodCode, methodCode)
}
func (c *Context) GetMethodCode() string {
	methodCode, exists := c.Get(CtxKeyMethodCode)
	if !exists {
		methodCode = c.Request.Header.Get(CtxKeyMethodCode)
	}
	return methodCode.(string)
}

func (c *Context) GetTraceId() string {
	traceId := c.Request.Header.Get(CtxKeyTraceID)
	if traceId == "" {
		traceId = GenerateUuid()
		c.Set(CtxKeyTraceID, traceId)
	}
	return traceId
}
