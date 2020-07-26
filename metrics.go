// Package httpmetrics 提供统计go程序http请求的指标
package httpmetrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// 指标的命名空间
	namespace = "http"
	subsystem = ""
)

var (
	requestTotalCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_total",
		Help:      "record requests total count",
	}, []string{"host"})
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_count",
		Help:      "record requests count",
	}, []string{"host", "method", "route", "code"})
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "requests_duration_seconds",
		Help:      "record requests duration in second",
		Buckets: []float64{
			.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 25, 30, 45, 60,
		},
	}, []string{"host", "method", "route", "code"})
)

// incRequestTotalCount 递增请求总数
func incRequestTotalCount(host string) {
	requestTotalCount.With(prometheus.Labels{"host": host}).Inc()
}

// incRequestCount 递增请求数量
func incRequestCount(host, method, route string, code int) {
	statusCode := strconv.Itoa(code)
	requestCount.With(prometheus.Labels{
		"host":   host,
		"method": method,
		"route":  route,
		"code":   statusCode,
	}).Inc()
}

// observeRequestDuration 设置请求持续时间
func observeRequestDuration(host, method, route string, code int, t time.Time) {
	d := float64(time.Since(t)) / float64(time.Second)
	statusCode := strconv.Itoa(code)
	requestDuration.With(prometheus.Labels{
		"host":   host,
		"method": method,
		"route":  route,
		"code":   statusCode,
	}).Observe(d)
}
