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
// 解析完整的路由，将 pattern 拆分成 parts []string
func parsePattern(pattern string) []string {
	// 将完整路由去掉"/"并拆分成切片
	items := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			// 如果第一位是'*'，则后面的part就不用看了
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	// 先将要添加的完整路由解析成 []string
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 获取请求方式为 method 的 Trie 树根节点；如果没有就创建一个根节点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 对请求方式为 method 的 Trie 树根节点进行插入，参数为：完整路由，解析路由，树高
	r.roots[method].insert(pattern, parts, 0)
	// 对请求 key 为 method + "-" + pattern 的 HandlerFunc 赋值
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	// 先将要搜索的完整路由解析成 []string
	searchParts := parsePattern(path)
	params := make(map[string]string)
	// 获取请求方式为 method 的 Trie 树根节点
	root, ok := r.roots[method]
	if !ok {
		// 如果获取不到，说明没有符合的路由，直接返回 nil
		return nil, nil
	}
	// 找到符合要获取的路由的节点
	n := root.search(searchParts, 0)
	if n != nil { // 找到了，进一步判断
		// 解析找到的节点本身的pattern，因为节点的pattern可以是带有:或者*的，而实际要搜索的path中，:或*会替换为相应的参数
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 如果part是参数匹配符:
			if part[0] == ':' {
				// 假如part是 ":name" ，则 params["name"] = 要搜索的完整路由解析成的 []string 的第 index 项
				// 即参数 :name 变成了amadeus
				params[part[1:]] = searchParts[index]
				// 处理完参数匹配符之后，还要继续处理后面的parts中剩余的部分
			}
			// 如果part是通配符*，且part不止一个
			if part[0] == '*' && len(part) > 1 {
				// 假如part是 "*filepath" ，则 params["filepath"] = 要搜索的完整路由解析成的 []string 的第 index+1 项及以后的项，并用/连接起来
				// 即参数 *filepath 变成了 go/amadeus.go
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				// 处理完通配符之后，就不需要再处理后面的parts中剩余的部分了
				break
			}
		}
		return n, params
	}
	return nil, nil // 没找到，直接返回 nil
}

func (r *router) handle(c *Context) {
	// 获取节点和参数
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
