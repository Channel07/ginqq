package main

import (
	gin "github.com/channel07/ginqq"
	"net/http"
)

func main() {
	r := gin.Default("Y122010101", "Y122")

	g := r.Group("/api")
	g.GET("/hello", func(c *gin.Context) {
		traceId := c.GetTraceId()
		http.Get("http://www.baidu.com")
		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!", "traceId": traceId})
	})
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
	})

	r.Run(":8080")

	//cfg := &gin.Config{
	//	ServiceCode: "Y122010101",
	//	PlatCode:    "Y122",
	//	MetricsConfig: &gin.MetricsConfig{
	//		Buckets: []float64{100, 200, 500, 1000, 2000, 3000, 5000, 10000},
	//		CustomLabels: map[string]string{
	//			"operate_type": "OPERATE_TYPE",
	//		},
	//	},
	//	HttpClientEnhanceConfig: &gin.HttpClientEnhanceConfig{
	//		Transport: http.DefaultTransport,
	//	},
	//}
	//r := gin.NewEngineWithConfig(cfg)
	//
	//r.GET("/hello", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
	//})
	//
	//r.Run(":8080")
}
