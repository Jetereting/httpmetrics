package httpmetrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Handle 包装http.Handler, 请求到达时http请求总数指标加一
// 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func Handle(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		incRequestTotalCount()
		method := r.Method
		path := r.URL.Path
		sw := &statusWriter{ResponseWriter: w}
		h.ServeHTTP(sw, r)
		traceRequest(strconv.Itoa(sw.status), method, path, t0)
	}
}

// GinMiddleware gin的中间件, 请求到达时http请求总数指标加一
// 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t0 := time.Now()
		incRequestTotalCount()
		method := c.Request.Method
		path := c.Request.URL.Path
		c.Next()
		traceRequest(strconv.Itoa(c.Writer.Status()), method, path, t0)
	}
}
