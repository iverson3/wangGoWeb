package wang

import "log"

// 跨域中间件
// 通过设置跨域相关的header头，统一解决跨域问题
func CORS() HandlerFunc {
	return func(c *Context) {
		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Credentials", "true")
		c.SetHeader("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.SetHeader("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Req.Method == "OPTIONS" {
			c.Fail(204, "cors error")
			log.Println("err3")
			return
		}

		c.Next()
	}
}
