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

var (
	requestTotalCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_total",
		Help:      "record requests total count",
	})
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_count",
		Help:      "record requests count",
	}, []string{"method", "route", "code"})
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_duration_seconds",
		Help:      "record requests duration in second",
		Buckets: []float64{
			.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 25, 30, 45, 60,
		},
	}, []string{"method", "route", "code"})
)

// incRequestTotalCount 递增请求总数
func incRequestTotalCount() {
	requestTotalCount.Inc()
}

// incRequestCount 递增请求数量
func incRequestCount(method, route, code string) {
	requestCount.With(prometheus.Labels{"method": method, "route": route, "code": code}).Inc()
}

// setRequestDuration 设置请求持续时间
func setRequestDuration(method, route, code string, d time.Duration) {
	v := float64(d) / float64(time.Second)
	requestDuration.With(prometheus.Labels{"method": method, "route": route, "code": code}).Observe(v)
}

func traceRequest(code, method, route string, ct time.Time) {
	incRequestCount(method, route, code)
	setRequestDuration(method, route, code, time.Since(ct))
}
