package main

import (
	gin "github.com/channel07/ginqq"
	"net/http"
)

func main() {
	r := gin.Default("A186010101", "channel07_ginqq")
	r.POST("/hello", gin.SetMethodCode("I00101"), Hello)
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
	//r.GET("/hello", gin.SetMethodCode("I00101"), Hello)
	//
	//r.Run(":8080")
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
}
