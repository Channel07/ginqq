package main

import (
	gin "github.com/channel07/ginqq"
	"net/http"
)

func main() {
	r := gin.Default("A186010101", "channel07-ginqq")

	r.POST("/hello", gin.SetMethodCode("I00101"), func(c *gin.Context) {
		//svcCode := r.Config.SvcCode
		//appName := strings.ReplaceAll(r.Config.AppName, "-", "_")
		//fmt.Printf("%s_%s\n", svcCode, appName)
		//fmt.Println(c.GetMethodCode())

		//url := c.Request.URL
		//fmt.Println(url)
		//fmt.Println(url.Host)
		//fmt.Println(url.Hostname())
		//fmt.Println(url.Query())
		//fmt.Println("..................")
		//
		//if err := c.Request.ParseForm(); err != nil {
		//	fmt.Println(err)
		//}
		//
		//form := make(map[string]interface{})
		//for key, values := range c.Request.Form {
		//	if len(values) == 1 {
		//		form[key] = values[0]
		//	} else {
		//		form[key] = values
		//	}
		//}
		//fmt.Println(form)
		//
		//var jsonData map[string]interface{}
		//if err := c.ShouldBindJSON(&jsonData); err != nil {
		//	fmt.Println(err)
		//	// not json data and form data: EOF
		//	// GET form: invalid character 'a' looking for beginning of value
		//	// POST form: EOF
		//}
		//fmt.Println(jsonData)
		//
		//xx, _ := json.Marshal(jsonData)
		//fmt.Println(string(xx))

		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
	})

	//g := r.Group("/api")
	//{
	//	g.GET("/hello", func(c *gin.Context) {
	//		traceId := c.GetTraceId()
	//		http.Get("http://www.baidu.com")
	//		c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!", "traceId": traceId})
	//	})
	//}

	r.Run(":80")

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
