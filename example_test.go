package httpmetrics_test

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gogs.xiaoyuanjijiehao.com/aag/httpmetrics"
)

func ExampleGinMiddleware() {
	r := gin.New()
	middleware := httpmetrics.GinMiddleware(&r.RouterGroup)
	r.Use(middleware)
	root := r.Group("/api")
	root.GET("/hello/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "hello %s", c.Param("id"))
	})
	r.Run(":8080")
}

func ExampleNewServeMux() {
	mux := httpmetrics.NewServeMux(http.DefaultServeMux)
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world!"))
	})
	http.ListenAndServe(":8080", mux)
}
