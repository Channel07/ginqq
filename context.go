// Package ginqq 自定义的context，用于传递系统上下文信息，如服务编码，traceId等
package ginqq

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
)

const (
	XMethodCode = "Method-Code"
	XTraceID    = "Transaction-ID"
)

type Context struct {
	*gin.Context
}

func Wrap(c *gin.Context) *Context {
	return &Context{Context: c}
}

// SetMethodCode 框架中间件数据传递
func (c *Context) SetMethodCode(methodCode string) {
	c.Set(XMethodCode, methodCode)
}
func (c *Context) GetMethodCode() string {
	methodCode, exists := c.Get(XMethodCode)
	if !exists {
		methodCode = c.Request.Header.Get(XMethodCode)
	}
	return methodCode.(string)
}

func (c *Context) GetTraceId() string {
	traceId := c.Request.Header.Get(XTraceID)
	if traceId == "" {
		traceId = uuid4()
		c.Set(XTraceID, traceId)
	}
	return traceId
}

// GetRawDataReusable get the request-body and reset it.
func (c *Context) GetRawDataReusable() ([]byte, error) {
	body, err := c.GetRawData()
	length := len(body)
	if length > 0 {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Request.ContentLength = int64(length)
		c.Request.Header.Set("Content-Length", strconv.Itoa(length))
	}
	return body, err
}
