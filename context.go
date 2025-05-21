// Package ginqq 自定义的context，用于传递系统上下文信息，如服务编码，traceId等
package ginqq

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"strings"
)

const (
	XTraceID         = "Trace-ID"
	XTransactionID   = "Transaction-ID"
	XFCode           = "User-Agent"
	XMethodCode      = "Method-Code"
	XMethodName      = "Method-Name"
	XResponsePayload = "Response-Payload"
)

type Context struct {
	*gin.Context
}

func Wrap(c *gin.Context) *Context {
	return &Context{Context: c}
}

func (c *Context) IndentedJSON(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.IndentedJSON(code, obj)
}

func (c *Context) SecureJSON(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.SecureJSON(code, obj)
}

func (c *Context) JSONP(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.JSONP(code, obj)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.JSON(code, obj)
}

func (c *Context) AsciiJSON(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.AsciiJSON(code, obj)
}

func (c *Context) PureJSON(code int, obj interface{}) {
	c.Set(XResponsePayload, obj)
	c.Context.PureJSON(code, obj)
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

func (c *Context) GetTraceID() string {
	traceID := c.GetString(XTraceID)
	if traceID == "" {
		traceID = c.GetHeader(XTraceID)
		if traceID == "" {
			traceID = uuid4()
		}
		c.Set(XTraceID, traceID)
	}
	return traceID
}

func (c *Context) GetTransactionID() string {
	transactionID := c.GetString(XTransactionID)
	if transactionID == "" {
		transactionID = c.GetHeader(XTransactionID)
		if transactionID == "" {
			transactionID = uuid4()
		}
		c.Set(XTransactionID, transactionID)
	}
	return transactionID
}

func (c *Context) GetFCode() string {
	return strings.ToUpper(c.GetHeader(XFCode))
}

func (c *Context) GetMethodCode() string {
	methodCode := c.GetString(XMethodCode)
	if methodCode == "" {
		methodCode = c.GetHeader(XMethodCode)
	}
	return methodCode
}

func (c *Context) GetMethodName() string {
	return c.GetString(XMethodName)
}

func (c *Context) GetResponsePayload() interface{} {
	responsePayload, _ := c.Get(XResponsePayload)
	return responsePayload
}
