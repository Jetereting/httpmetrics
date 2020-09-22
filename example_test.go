package httpmetrics_test

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gogs.xiaoyuanjijiehao.com/aag/httpmetrics"
)

func ExampleGinMiddleware() {
	r := gin.New()
	ips := []string{
		"127.0.0.1",
	}
	middleware := httpmetrics.GinMiddleware(&r.RouterGroup,
		httpmetrics.AllowIPsOption(ips),
		httpmetrics.EnableAllowListOption(true),
	)
	// 初始化中间件时可以自定义抓取路径
	r.Use(middleware)
	root := r.Group("/api")
	root.GET("/hello/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "hello %s", c.Param("id"))
	})
	r.Run(":8080")
}

func ExampleNewServeMux() {
	ips := []string{
		"127.0.0.1",
	}
	mux := httpmetrics.NewServeMux(http.NewServeMux(),
		httpmetrics.AllowIPsOption(ips),
		httpmetrics.EnableAllowListOption(true),
	)
	http.ListenAndServe(":8080", mux)
}
