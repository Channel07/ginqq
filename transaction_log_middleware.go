package ginqq

// TransactionLogMiddleware 流水日志中间件
func transactionLogMiddleware() func(*Context) {
	return func(c *Context) {
		before(c)
		c.Next()
		after(c)
	}
}

func before(c *Context) {
	println("Transaction log before")
}

func after(c *Context) {
	println("Transaction log after")
}
