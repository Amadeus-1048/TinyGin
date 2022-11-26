package tinyGin

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by tinyGin
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine implement the interface of ServeHTTP
// 定义了一个结构体Engine，实现了ServeHTTP接口
type Engine struct {
	router map[string]HandlerFunc
}

// New is the constructor of tinyGin.Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

// Run defines the method to start a http server
// Engine 必须实现 ServeHTTP方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
