package main

import (
	"log"
	"net/http"
	"time"
	"tinyGin"
)

/*
Handler的参数变成成了gee.Context，提供了查询Query/PostForm参数的功能。
gee.Context封装了HTML/String/JSON函数，能够快速构造HTTP响应。
*/

func middlewareForV2() tinyGin.HandlerFunc {
	return func(c *tinyGin.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := tinyGin.Default()

	r.GET("/", func(c *tinyGin.Context) {
		c.String(http.StatusOK, "Hello TinyGin\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *tinyGin.Context) {
		names := []string{"Amadeus"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}

/*
测试方法
curl -i http://localhost:9999/

curl "http://localhost:9999/hello?name=amadeus"

curl "http://localhost:9999/login" -X POST -d 'username=amadeus&password=1240'

curl "http://localhost:9999/xxx"
*/
