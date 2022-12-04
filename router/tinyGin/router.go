package tinyGin

import (
	"net/http"
	"strings"
)

/*
将和路由相关的方法和结构提取了出来，放到了一个新的文件中router.go，
方便我们下一次对 router 的功能进行增强，例如提供动态路由的支持。
router 的 handle 方法作了一个细微的调整，即 handler 的参数，变成了 Context。
*/

type router struct {
	roots    map[string]*node       // 存储每种请求方式的 Trie 树根节点
	handlers map[string]HandlerFunc // 存储每种请求方式的 HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + c.Path
		// 在调用匹配到的handler前，将解析出来的路由参数赋值给了c.Params
		c.Params = params
		// 这样就能够在handler中，通过Context对象访问到具体的值了
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
