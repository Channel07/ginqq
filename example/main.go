package main

import (
	gin "github.com/channel07/ginqq"
	"net/http"
)

func main() {
	r := gin.Default("A186010101", "channel07_ginqq")
	r.POST("/hello", gin.MethodCode("I00101"), Hello)
	r.Run(":8080")

}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from ginqq!"})
}
