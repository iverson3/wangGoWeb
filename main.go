package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"wang"
)

type student struct {
	Name string
	Age int8
}

// 自定义的模板渲染函数 - 格式化时间
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func onlyForV2() wang.HandlerFunc {
	return func(c *wang.Context) {
		t := time.Now()

		c.Fail(500, "Internal Server Error")

		log.Printf("[%d] %s in %v for group v2\n", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	//engine := wang.New()
	// global middleware
	//engine.Use(wang.Logger())   // 日志中间件
	//engine.Use(wang.Recovery()) // 错误处理中间件

	engine := wang.Default()

	//engine.Use(wang.CORS())

	// 设置自定义的模板渲染函数
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载templates目录下的所有模板文件
	engine.LoadHTMLGlob("templates/*")
	// 开启对静态资源文件的处理
	engine.Static("/assets", "./static")  // 相对路径
	//engine.Static("/assets", "/usr/geektutu/blog/static")  // 绝对路径

	stu1 := &student{Name: "tom", Age: 22}
	stu2 := &student{Name: "stefan", Age: 27}

	engine.GET("/", func(c *wang.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	engine.GET("/students", func(c *wang.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", wang.H{
			"title": "wang",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	engine.GET("/date", func(c *wang.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", wang.H{
			"title": "wang",
			"now": time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})


	engine.GET("/panic", func(c *wang.Context) {
		names := []string{"xxx"}
		_, _ = c.String(http.StatusOK, names[100])
	})

	v1 := engine.Group("/v1")
	{
		//v1.GET("/", func(c *wang.Context) {
		//	c.HTML(http.StatusOK, "<h1>Hello wang.</h1>")
		//})
		v1.GET("/hello", func(c *wang.Context) {
			// expect /hello?name=tom
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := engine.Group("/v2")
	v2.Use(onlyForV2())  // middleware for v2 group
	{
		v2.GET("/hello/:name", func(c *wang.Context) {
			// expect /hello/jack
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *wang.Context) {
			c.JSON(http.StatusOK, wang.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	// 静态资源文件的处理已经交由 engine.Static()来处理了；所以这里不需要了
	//engine.GET("/assets/*filepath", func(c *wang.Context) {
	//	c.JSON(http.StatusOK, wang.H{
	//		"filepath": c.Param("filepath"),
	//	})
	//})

	err := engine.Run(":9999")
	if err != nil {
		log.Fatal(err)
	}
}




















