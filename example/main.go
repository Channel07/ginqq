package main

import (
	"fmt"
	gin "github.com/channel07/ginqq"
	"net/http"
	"strings"
)

func main() {
	r := gin.Default("A186010101", "channel07-ginqq")

	r.GET("/hello", gin.SetMethodCode("I00101"), func(c *gin.Context) {
		svcCode := r.Config.SvcCode
		appName := strings.ReplaceAll(r.Config.AppName, "-", "_")
		fmt.Printf("%s_%s\n", svcCode, appName)
		fmt.Println(c.GetMethodCode())
		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
	})

	g := r.Group("/api")
	{
		g.GET("/hello", func(c *gin.Context) {
			traceId := c.GetTraceId()
			http.Get("http://www.baidu.com")
			c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!", "traceId": traceId})
		})
	}

	r.Run(":8080")

	//cfg := &gin.Config{
	//	SvcCode: "Y122010101",
	//	AppName: "channel07-ginqq",
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
