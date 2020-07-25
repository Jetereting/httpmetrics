package httpmetrics

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHTTPHandleFuncMetrics(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	ms := httptest.NewServer(promhttp.Handler())
	defer ms.Close()
	h := func(w http.ResponseWriter, r *http.Request) {
		n := rand.Int63n(50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		if n%2 == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`hello world!`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
	ts := httptest.NewServer(Handle(http.HandlerFunc(h)))
	for i := 0; i < 20; i++ {
		resp, _ := http.Post(ts.URL, "application/json", nil)
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	ts.Close()
	//
	resp, err := http.Get(ms.URL)
	if err != nil {
		t.Errorf("http.Get() error = %v", err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Errorf("ioutil.ReadAll() error = %v", err)
		return
	}
	t.Logf("%s", b)
}

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

func TestGinMiddlewareMetrics(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	r := gin.New()
	r.GET("/metrics", func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	})
	r.Use(GinMiddleware())
	r.GET("/test", func(c *gin.Context) {
		n := rand.Int63n(50)
		time.Sleep(time.Duration(n) * time.Millisecond)
		if n%2 == 0 {
			c.String(http.StatusOK, "%s", "hello world!")
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
	for i := 0; i < 20; i++ {
		performRequest(r, "GET", "/test")
	}
	//
	w := performRequest(r, "GET", "/metrics")
	t.Logf("statusCode %v", w.Code)
	t.Logf("body: %v", w.Body.String())
}
