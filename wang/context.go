package wang

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 上下文
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        r,
		Path:       r.URL.Path,
		Method:     r.Method,
	}
}

// 獲取post參數
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 獲取url參數
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 獲取url路徑參數 比如 /v/1.0.1/id/12138
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 為response header設置status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 設置response header
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}


// 以下為支持的四種response數據格式
func (c *Context) String(code int, format string, values ...interface{}) (int, error) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	return c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) Data(code int, data []byte) (int, error) {
	c.Status(code)
	return c.Writer.Write(data)
}
func (c *Context) HTML(code int, html string) (int, error) {
	//正确的调用顺序应该是Header().Set 然后WriteHeader() 最后是Write()
	//在 WriteHeader() 后调用 Header().Set 是不会生效的
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	return c.Writer.Write([]byte(html))
}