package httpmetrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// ServeMux http 路由
type ServeMux struct {
	metrics string
	*http.ServeMux
}

func checkMetricsPath(s ...string) string {
	p := "/metrics"
	if len(s) > 0 && s[0] != "" {
		p = s[0]
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
	}
	return p
}

// NewServeMux 创建新的ServeMux, 请求到达时http请求总数指标加一
// 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func NewServeMux(mux *http.ServeMux, metricsPath ...string) *ServeMux {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	p := checkMetricsPath(metricsPath...)
	out := &ServeMux{metrics: p, ServeMux: mux}
	out.Handle(p, promhttp.Handler())
	return out
}

// ServeHTTP 处理http请求, 并设置相关指标
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next, pattern := mux.Handler(r)
	if pattern == mux.metrics {
		next.ServeHTTP(w, r)
		return
	}
	t0 := time.Now()
	incRequestTotalCount()
	method := r.Method
	sw := &statusWriter{ResponseWriter: w}
	next.ServeHTTP(sw, r)
	traceRequest(strconv.Itoa(sw.status), method, pattern, t0)
}

// GinMiddleware gin的中间件, 注册采集器指标接口
// 在请求到达时http请求总数指标加一, 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func GinMiddleware(r *gin.RouterGroup, metricsPath ...string) gin.HandlerFunc {
	p := checkMetricsPath(metricsPath...)
	h := promhttp.Handler()
	r.GET(p, func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == p {
			c.Abort()
			return
		}
		t0 := time.Now()
		incRequestTotalCount()
		method := c.Request.Method
		c.Next()
		traceRequest(strconv.Itoa(c.Writer.Status()), method, path, t0)
	}
}
