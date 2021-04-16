package wang

import (
	"net/http"
)

// 請求處理方法定義
type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: NewRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	// 第二个参数代表处理所有的HTTP请求的实例，如果是nil則代表使用标准库中的实例处理
	return http.ListenAndServe(addr, engine)
}

// 實現http的Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	c := NewContext(w, req)
	engine.router.Handle(c)
}