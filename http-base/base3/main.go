package main

import (
	"fmt"
	"net/http"
	"tinyGin"
)

/*
至此实现了路由映射表，
提供了用户注册静态路由的方法，
包装了启动服务的函数
*/

func main() {
	r := tinyGin.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "Path of URL = %q\n", req.URL.Path)
		if err != nil {
			return
		}
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for i, v := range req.Header {
			_, err := fmt.Fprintf(w, "Header[%q] = %q\n", i, v)
			if err != nil {
				continue
			}
		}
	})

	r.Run(":9999")
}
