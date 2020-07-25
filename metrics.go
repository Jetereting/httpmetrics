// Package httpmetrics 提供统计go程序http请求的指标
package httpmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// 指标的命名空间
	namespace = "ant"
	subsystem = "http"
)

// http请求统计的标签
var labels = []string{"method", "path", "status"}

var (
	requestTotalCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_total_count",
		Help:      "record request total count",
	})
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_count",
		Help:      "record request count",
	}, []string{"method", "path", "status"})
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_duration_seconds",
		Help:      "record request duration in second",
		Buckets: []float64{
			// 默认的Buckets
			.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 25, 30, 45, 60,
		},
	}, []string{"method", "path", "status"})
)

// incRequestTotalCount 递增请求总数
func incRequestTotalCount() {
	requestTotalCount.Inc()
}

// incRequestCount 递增请求数量
func incRequestCount(method, path, status string) {
	requestCount.With(prometheus.Labels{"method": method, "path": path, "status": status}).Inc()
}

// setRequestDuration 设置请求持续时间
func setRequestDuration(method, path, status string, d time.Duration) {
	v := float64(d) / float64(time.Second)
	requestDuration.With(prometheus.Labels{"method": method, "path": path, "status": status}).Observe(v)
}

func traceRequest(status, method, path string, ct time.Time) {
	incRequestCount(method, path, status)
	setRequestDuration(method, path, status, time.Since(ct))
}
