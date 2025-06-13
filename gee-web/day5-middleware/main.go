package main

/*
(1) global middleware Logger
$ curl http://localhost:9999/
<h1>Hello Gee</h1>

>>> log
2019/08/17 01:37:38 [200] / in 3.14µs
*/

/*
(2) global + group middleware
$ curl http://localhost:9999/v2/hello/geektutu
{"message":"Internal Server Error"}

>>> log
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 61.467µs for group v2
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 281µs
*/

import (
	"log"
	"net/http"
	"time"

	"gee"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	// 创建Gee引擎实例
	r := gee.New()

	// 注册全局中间件（对所有路由生效）
	r.Use(gee.Logger())       // 请求日志记录中间件
	r.Use(gee.MyMiddleware()) // 自定义全局中间件

	// 注册根路由
	r.GET("/", func(c *gee.Context) {
		// 返回HTML响应，状态码200
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	// 创建v2版本路由组（路径前缀为/v2）
	v2 := r.Group("/v2")
	// 注册v2组专属中间件（仅对/v2开头的路由生效）
	v2.Use(onlyForV2())

	// 在v2路由组中注册子路由
	{
		// 注册动态路由，匹配格式为 /v2/hello/{name}
		v2.GET("/hello/:name", func(c *gee.Context) {
			// 从URL参数中获取name值，示例：/hello/geektutu -> name=geektutu
			// 返回格式化字符串响应，状态码200
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	// 启动服务器并监听9999端口
	r.Run(":9999")
}
