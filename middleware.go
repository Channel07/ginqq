package ginqq

// LoggingServerMiddleware 系统日志中间件
func LoggingServerMiddleware() func(*Context) {
	return func(ctx *Context) {
		// 记录日志
		println("InJournalLogMiddleware")
		println(globalConfig.ServiceCode)
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
