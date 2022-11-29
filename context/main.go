package main

import (
	"net/http"
	"tinyGin"
)

/*
Handler的参数变成成了gee.Context，提供了查询Query/PostForm参数的功能。
gee.Context封装了HTML/String/JSON函数，能够快速构造HTTP响应。
*/

func main() {
	r := tinyGin.New()
	r.GET("/", func(c *tinyGin.Context) {
		c.HTML(http.StatusOK, "<h1>Hello TinyGin</h1>")
	})

	r.GET("/hello", func(c *tinyGin.Context) {
		// expect /hello?name=amadeus
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *tinyGin.Context) {
		c.Json(http.StatusOK, tinyGin.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
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
