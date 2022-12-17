package tinyGin

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by tinyGin
// 定义路由映射的处理方法
// 对Web服务来说，就是根据请求*http.Request，构造响应http.ResponseWriter
type HandlerFunc func(ctx *Context)

/*
分组控制(Group Control)是 Web 框架应提供的基础功能之一。所谓分组，是指路由的分组
想要实现路由分组功能，首先需要抽象出一个RouterGroup的类型出来，这个类型应该能满足如下几个功能：
	正确的分组
	存储group的相关信息
	满足多层分组的需要（可以在group里面再新建group）
	能够对某一类group加中间件进行处理
*/

// RouterGroup 代表分组类型，包含四部分信息
type RouterGroup struct {
	prefix      string        // 当前group的前缀
	middlewares []HandlerFunc // support middleware	 当前分组需要执行的中间件处理函数
	parent      *RouterGroup  // support nesting	当前group的父group
	engine      *Engine       // all groups share an Engine instance	所有的group共享一个Engine对象，为了便于操作，可以直接在group中获取engine中的信息
}

// Engine implement the interface of ServeHTTP
// 定义了一个结构体Engine，实现了ServeHTTP接口
// Engine现在作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力
type Engine struct {
	router *router // 路由映射表
	// 将group相关的信息也加入到engine中，这里需要注意的时候，一个engine就相当于一个没有前缀的分组
	*RouterGroup                     // 所以在engine中也支持group相关的所有方法  这里用到了结构体的继承
	groups        []*RouterGroup     // store all groups  要有一个groups存放所有的group信息
	htmlTemplates *template.Template // for html render	将所有的模板加载进内存
	funcMap       template.FuncMap   // for html render	是所有的自定义模板渲染函数
}

// New is the constructor of tinyGin.Engine
func New() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{
		engine: e,
	}
	e.groups = []*RouterGroup{
		e.RouterGroup,
	}
	return e
}

// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	e := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: e,
	}
	e.groups = append(e.groups, newGroup)
	return newGroup
}

// 添加路由
// 新的addRoute函数，调用了group.engine.router.addRoute来实现了路由的映射
// 由于Engine从某种意义上继承了RouterGroup的所有属性和方法，因为 (*Engine).engine 是指向自己的。
// 这样实现，我们既可以像原来一样添加路由，也可以通过分组添加路由。
func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
// 当用户调用(*Engine).GET()方法时，会将路由和处理方法注册到映射表 router 中
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use is defined to add middleware to the group
// 中间件应该与Group对象绑定，因为需要中间件的时候，肯定是要对一类路由进行处理。如果仅仅单个路由需要，那完全可以将逻辑放入到对应路由的处理函数里面
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	// 将中间件添加到group中
	group.middlewares = append(group.middlewares, middlewares...)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// Engine实现的 ServeHTTP 方法的作用：解析请求的路径，查找路由映射表，如果查到，就执行注册的处理方法。如果查不到，就返回 404 NOT FOUND
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewares := []HandlerFunc{}
	// 查出本次请求所有需要调用的中间件
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			// 保持下来放到context中
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	// 在 Context 中添加了成员变量 engine *Engine，这样就能够通过 Context 访问 Engine 中的 HTML 模板。实例化 Context 时，还需要给 c.engine 赋值
	c.engine = e
	// 查出本次请求对应的处理函数，然后再依次开始请求
	e.router.handle(c)
}

// Run defines the method to start a http server
// (*Engine).Run()方法，是 ListenAndServe 的包装
// Engine 必须实现 ServeHTTP方法，才能使用http.ListenAndServe
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// SetFuncMap 设置自定义渲染函数
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// LoadHTMLGlob 加载模板
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
