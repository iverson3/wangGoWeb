package main

import (
	"log"
	"net/http"
	"wang"
)

func main() {
	engine := wang.New()

	engine.GET("/", func(c *wang.Context) {
		c.HTML(http.StatusOK, "<h1>hello Wang</h1>")
	})
	engine.GET("/hello", func(c *wang.Context) {
		// expect /hello?name=tom
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	engine.GET("/hello/:name", func(c *wang.Context) {
		// expect /hello/jack
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	engine.GET("/assets/*filepath", func(c *wang.Context) {
		c.JSON(http.StatusOK, wang.H{
			"filepath": c.Param("filepath"),
		})
	})

	err := engine.Run(":9999")
	if err != nil {
		log.Fatal(err)
	}
}




















