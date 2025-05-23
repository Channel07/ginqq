package ginqq

import "strings"

// LoggingServerMiddleware 系统日志中间件。
func LoggingServerMiddleware() func(*Context) {
	return func(ctx *Context) {
		// 记录日志
		println("InJournalLogMiddleware")
		println(cnf.SvcCode)
		ctx.Next()
	}
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

// SetMethodCode 是一个中间件，用于设置接口编码。
func SetMethodCode(methodCode string) func(*Context) {
	return func(c *Context) {
		c.Set(XMethodCode, strings.ToUpper(methodCode))
		c.Next()
	}
}
