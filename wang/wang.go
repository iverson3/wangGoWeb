package wang

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// 請求處理方法定義
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string               // 前缀
	middlewares []HandlerFunc   // 支持中间件
	parent *RouterGroup         // 支持嵌套
	engine *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup  // store all groups

	htmlTemplates *template.Template  // for html render  将所有的模板加载进内存
	funcMap        template.FuncMap   // for html render  所有的自定义模板渲染函数
}

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 调用Default()会默认为engine注册Logger 和 Recovery 这两个中间件
func Default() *Engine {
	engine := New()
	// 默认使用 Logger 和 Recovery 中间件
	engine.Use(Logger(), Recovery())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
// serve static file
// Static方法是暴露给用户的。用户可以将磁盘上的某个文件夹root映射到路由relativePath
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")

	// 所有静态资源文件的请求处理
	group.GET(urlPattern, handler)
}

// 设置自定义的模板渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}
// 加载所有符合pattern的html模板； 参数是模板文件的路径或者目录 比如 "templates/*"
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) Run(addr string) error {
	// 第二个参数代表处理所有的HTTP请求的实例，如果是nil則代表使用标准库中的实例处理
	return http.ListenAndServe(addr, engine)
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 實現http的Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	var middlewares []HandlerFunc
	// 循环判断当前请求适用于哪些中间件(属于哪些路由组)
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := NewContext(w, req)
	c.middlewareHandlers = middlewares
	c.engine = engine
	engine.router.Handle(c)
}