package httpmetrics

import (
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gogs.xiaoyuanjijiehao.com/aag/httpmetrics/internal/ip"
)

// Options 选项
type Options struct {
	metricsPath  string
	allowIPList  []string
	enableIPList bool
}

// mergeOptions 创建选项
func mergeOptions(opts ...Option) *Options {
	o := &Options{
		metricsPath: "/metrics",
	}
	for _, optFn := range opts {
		optFn(o)
	}
	return o
}

// Option 选项参数
type Option func(opts *Options)

// MetricPathOption 抓取路径
func MetricPathOption(metricsPath string) Option {
	return func(opts *Options) {
		s := strings.TrimSpace(metricsPath)
		if s == "" {
			return
		}
		s = filepath.Join("/", s)
		opts.metricsPath = s
	}
}

// AllowIPsOption 允许ip列表
func AllowIPsOption(ips []string) Option {
	return func(opts *Options) {
		opts.allowIPList = ips
	}
}

// EnableAllowListOption 启用允许列表
func EnableAllowListOption(v bool) Option {
	return func(opts *Options) {
		opts.enableIPList = v
	}
}

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
	*http.ServeMux

	metricsPath  string
	allowIPs     *ip.AllowList
	enableIPList bool
}

// NewServeMux 创建新的ServeMux, 请求到达时http请求总数指标加一
// 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func NewServeMux(mux *http.ServeMux, opts ...Option) *ServeMux {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	options := mergeOptions(opts...)
	out := &ServeMux{
		metricsPath:  options.metricsPath,
		allowIPs:     ip.LoadArray(options.allowIPList),
		enableIPList: options.enableIPList,
		ServeMux:     mux,
	}
	out.Handle(options.metricsPath, promhttp.Handler())
	return out
}

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ServeHTTP 处理http请求, 并设置相关指标
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next, pattern := m.Handler(r)
	if pattern == m.metricsPath {
		if m.enableIPList &&
			!m.allowIPs.ContainsString(getClientIP(r)) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
		return
	}
	t0 := time.Now()
	host := r.Host
	incRequestTotalCount(host)
	method := r.Method
	sw := &statusWriter{ResponseWriter: w}
	next.ServeHTTP(sw, r)
	code := sw.status
	incRequestCount(host, method, pattern, code)
	if code == http.StatusNotFound {
		return
	}
	observeRequestDuration(host, method, pattern, code, t0)
}

// GinMiddleware gin的中间件, 注册采集器指标接口
// 在请求到达时http请求总数指标加一, 返回响应后, 请求数量指标加一, 且统计http请求持续时间
func GinMiddleware(r *gin.RouterGroup, opts ...Option) gin.HandlerFunc {
	options := mergeOptions(opts...)
	h := promhttp.Handler()
	allowIPList := ip.LoadArray(options.allowIPList)
	r.GET(options.metricsPath, func(c *gin.Context) {
		if options.enableIPList && !allowIPList.ContainsString(c.ClientIP()) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		h.ServeHTTP(c.Writer, c.Request)
	})
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == options.metricsPath {
			c.Abort()
			return
		}
		t0 := time.Now()
		host := c.Request.Host
		incRequestTotalCount(host)
		method := c.Request.Method
		c.Next()
		code := c.Writer.Status()
		incRequestCount(host, method, path, code)
		if code == http.StatusNotFound {
			return
		}
		observeRequestDuration(host, method, path, code, t0)
	}
}

// BeegoMiddleware 创建beego的中间件
func BeegoMiddleware(opts ...Option) func(h http.Handler) func(http.ResponseWriter, *http.Request) {
	options := mergeOptions(opts...)
	ph := promhttp.Handler()
	allowIPList := ip.LoadArray(options.allowIPList)
	return func(h http.Handler) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == options.metricsPath {
				if allowIPList.ContainsString(getClientIP(r)) {
					ph.ServeHTTP(w, r)
					return
				}
				http.NotFound(w, r)
				return
			}
			t0 := time.Now()
			host := r.Host
			incRequestTotalCount(host)
			method := r.Method
			sw := &statusWriter{ResponseWriter: w}
			h.ServeHTTP(w, r)
			code := sw.status
			incRequestCount(host, method, path, code)
			if code == http.StatusNotFound {
				return
			}
			observeRequestDuration(host, method, path, code, t0)
		}
	}

}
