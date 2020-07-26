package httpmetrics

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestServerMuxMetrics(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()
	s := NewServeMux(mux)
	s.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		n := rand.Int63n(50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		if n%2 == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`hello world!`))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
	})
	for i := 0; i < 20; i++ {
		performRequest(s, "GET", "/test")
	}
	//
	w := performRequest(s, "GET", "/metrics")
	t.Logf("statusCode %v", w.Code)
	t.Logf("body: %v", w.Body.String())
}

func TestGinMiddlewareMetrics(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	r := gin.New()
	metrics := GinMiddleware(&r.RouterGroup)
	r.Use(metrics)
	r.GET("/test/:id", func(c *gin.Context) {
		n := rand.Int63n(50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		if n%2 == 0 {
			c.String(http.StatusOK, "%d: hello world!", c.Param("id"))
			return
		}
		c.AbortWithStatus(http.StatusBadGateway)
	})
	for i := 0; i < 20; i++ {
		performRequest(r, "GET", "/test/"+strconv.Itoa(i))
	}
	//
	w := performRequest(r, "GET", "/metrics")
	t.Logf("statusCode %v", w.Code)
	t.Logf("body: %v", w.Body.String())
}
