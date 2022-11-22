package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 设置2个路由，/和/hello，分别绑定 indexHandler 和 helloHandler ， 根据不同的HTTP请求会调用不同的处理函数。
	http.HandleFunc("/", indexHandler)      // 访问/，响应是Path of URL = /
	http.HandleFunc("/hello", helloHandler) // /hello的响应则是请求头(header)中的键值对信息
	// 下面一行是用来启动 Web 服务的
	// 第一个参数是地址，:9999表示在 9999 端口监听。
	// 第二个参数则代表处理所有的HTTP请求的实例，nil 代表使用标准库中的实例处理，即基于net/http标准库实现Web框架的入口。
	log.Fatal(http.ListenAndServe(":9999", nil))
	// 第二个参数handler的类型是一个接口，需要实现方法 ServeHTTP
	// 只要传入任何实现了 ServerHTTP 接口的实例，所有的HTTP请求，就都交给了该实例处理
	// 下面将在base2/main.go中实现一个handler实例engine
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "Path of URL = %q\n", req.URL.Path)
	if err != nil {
		return
	}
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for i, v := range req.Header {
		_, err := fmt.Fprintf(w, "Header[%q] = %q\n", i, v)
		if err != nil {
			continue
		}
	}
}

/*
测试方法
➜  TinyGin git:(main) ✗ curl http://localhost:9999/
Path of URL = "/"
➜  TinyGin git:(main) ✗ curl http://localhost:9999/hello
Header["User-Agent"] = ["curl/7.79.1"]
Header["Accept"] = ["*/ /*"]   这里多加了一个斜杠，防止注释结束
 */
