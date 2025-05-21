package ginqq

import "strings"

// MethodCode 是一个中间件，用于设置接口编码。
func MethodCode(I string) func(*Context) {
	return func(c *Context) {
		c.Set(XMethodCode, strings.ToUpper(strings.TrimSpace(I)))
		c.Next()
	}
}

// TransactionLogMiddleware 流水日志中间件。
func TransactionLogMiddleware() func(*Context) {
	return DispatchTransactionLog
}

func MetricsServerMiddleware() func(*Context) {
	return func(ctx *Context) {
		println("MetricsServerMiddleware")
		ctx.Next()
	}
}

func TracingServerMiddleware() func(*Context) {
	return func(ctx *Context) {
		println("TracingServerMiddleware")
		ctx.Next()
	}
}

func ApiStandardServerMiddleware() func(*Context) {
	return func(ctx *Context) {
		println("ApiStandardServerMiddleware")
		ctx.Next()
	}
}

func SecurityMiddleware() func(*Context) {
	return func(ctx *Context) {
		println("SecurityMiddleware")
		ctx.Next()
	}
}
