package httpmetrics_test

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gogs.xiaoyuanjijiehao.com/ops/httpmetrics"
)

func ExampleGinMiddleware() {
	r := gin.New()
	r.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})
	r.Use(httpmetrics.GinMiddleware())
	root := r.Group("/api")
	root.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.Run()
}

func ExampleHandle() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", nil)
}
