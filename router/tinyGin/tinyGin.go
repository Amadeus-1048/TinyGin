package tinyGin

import "net/http"

// HandlerFunc defines the request handler used by tinyGin
// 定义路由映射的处理方法
// 对Web服务来说，就是根据请求*http.Request，构造响应http.ResponseWriter
type HandlerFunc func(ctx *Context)

// Engine implement the interface of ServeHTTP
// 定义了一个结构体Engine，实现了ServeHTTP接口
type Engine struct {
	router *router // 路由映射表
}

// New is the constructor of tinyGin.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加路由
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
// 当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Engine实现的 ServeHTTP 方法的作用：解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

// Run defines the method to start a http server
// (*Engine).Run()方法，是 ListenAndServe 的包装
// Engine 必须实现 ServeHTTP方法，才能使用http.ListenAndServe
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
