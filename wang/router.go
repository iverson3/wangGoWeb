package wang

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// 通用的404错误处理handler
var Err404NotFoundFunc HandlerFunc = func(c *Context) {
	c.HTML(http.StatusNotFound, "404.html", H{
		"path": c.Path,
	})
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func NewRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	arr := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range arr {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %s-%s", method, pattern)
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
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

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) Handle(c *Context) {
	var _handler HandlerFunc

	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		_handler = Err404NotFoundFunc
	} else {
		c.Params = params
		key := c.Method + "-" + n.pattern

		// 将从路由中匹配得到的 Handler 添加到 c.middlewareHandlers列表中
		// 相当于把请求的处理逻辑也当作一个中间件，且位于中间件列表的最后一个
		if handler, ok := r.handlers[key]; ok {
			_handler = handler
		} else {
			_handler = Err404NotFoundFunc
		}
	}

	c.middlewareHandlers = append(c.middlewareHandlers, _handler)
	c.Next()
}