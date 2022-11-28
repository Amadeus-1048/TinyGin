package tinyGin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Context的结构体，结构体中包含三类元素。
首先是origin object（http.ResponseWriter、*http.Request），在之前我们已经知道这是一个route的处理函数所必须的输入参数；
然后是跟请求有关的信息request info，Path和Method都是从http.ResponseWriter取出的信息；
最后是跟响应有关的信息response info,StatusCode即响应码
*/

// H
// 给map[string]interface{}起了一个别名tinyGin.H，构建JSON数据时，显得更简洁
type H map[string]interface{}

// Context
// 要构造一个完整的响应，需要考虑消息头(Header)和消息体(Body)
// 而 Header 包含了状态码(StatusCode)，消息类型(ContentType)等几乎每次请求都需要设置的信息，需要进行有效的封装
// 针对使用场景，封装*http.Request和http.ResponseWriter的方法，简化相关接口的调用，只是设计 Context 的原因之一
// Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 为了简化接口，封装了一些http.Request方法以供使用

// PostForm 访问Query和PostForm参数的方法
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 访问Query和PostForm参数的方法
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 封装一些http.ResponseWriter方法使用，为了方便对于JSON、HTML等返回类型的支持，这些返回类型都是非常常见的，因此封装起来，减少调用的代码量

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// 快速构造String/Data/JSON/HTML响应的方法
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Json 快速构造String/Data/JSON/HTML响应的方法
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 快速构造String/Data/JSON/HTML响应的方法
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 快速构造String/Data/JSON/HTML响应的方法
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
