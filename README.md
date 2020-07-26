# httpmetrics

Prometheus 的 http 统计指标和中间件

采集器的抓取路径默认注册在路由组的下的`/metrics`路径, 如: `http://example.com/metrics`, 可以在初始化中间件时修改.

采集器在`Prometheus`的默认指标的基础上增加了三个指标, 使我们可以通过监控观察 HTTP 请求的主要指标, 如: 请求数量， 请求成功(失败)数量, 请求持续时间等.
借助`Prometheus`的查询语句可以分析主机、接口以及响应码等数据.

[指标示例](example.prom)

#### HTTP 指标

- `# TYPE http_requests_total counter` HTTP 请求总数量计数器

```
  http_requests_total{host="example.com"} 20
```

- `# TYPE http_requests_count counter` HTTP 请求数量计数器

```
  http_requests_count{code="200",host="example.com",method="GET",route="/test/:id"} 12
  http_requests_count{code="502",host="example.com",method="GET",route="/test/:id"} 8
```

- `# TYPE http_requests_duration_seconds histogram` HTTP 请求持续时间直方图

```
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.005"} 0
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.01"} 1
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.025"} 7
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.05"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.1"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.25"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="0.5"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="1"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="2.5"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="5"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="10"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="15"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="20"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="25"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="30"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="45"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="60"} 12
  http_requests_duration_seconds_bucket{code="200",host="example.com",method="GET",route="/test/:id",le="+Inf"} 12
  http_requests_duration_seconds_sum{code="200",host="example.com",method="GET",route="/test/:id"} 0.2993984
  http_requests_duration_seconds_count{code="200",host="example.com",method="GET",route="/test/:id"} 12
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.005"} 0
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.01"} 1
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.025"} 4
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.05"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.1"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.25"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="0.5"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="1"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="2.5"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="5"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="10"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="15"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="20"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="25"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="30"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="45"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="60"} 8
  http_requests_duration_seconds_bucket{code="502",host="example.com",method="GET",route="/test/:id",le="+Inf"} 8
  http_requests_duration_seconds_sum{code="502",host="example.com",method="GET",route="/test/:id"} 0.20078440000000003
  http_requests_duration_seconds_count{code="502",host="example.com",method="GET",route="/test/:id"} 8
```

#### Gin 中间件

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gogs.xiaoyuanjijiehao.com/ops/httpmetrics"
)

func main() {
  r := gin.New()
  middleware := httpmetrics.GinMiddleware(&r.RouterGroup)
  // 初始化中间件时可以自定义抓取路径
  // middleware := httpmetrics.GinMiddleware(&r.RouterGroup, "/monitor/metrics")
  r.Use(middleware)
  root := r.Group("/api")
  root.GET("/hello/:id", func(c *gin.Context) {
    c.String(http.StatusOK, "hello %s", c.Param("id"))
  })
  r.Run(":8080")
}
```

#### HTTP 中间件

```go
package main

import (
    "net/http"

    "gogs.xiaoyuanjijiehao.com/ops/httpmetrics"
)

func main() {
  mux := httpmetrics.NewServeMux(http.DefaultServeMux)
  // 可自定义抓取路径
  // mux := httpmetrics.NewServeMux(http.DefaultServeMux, "/monitor/metrics")
  mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("hello world!"))
  })
  http.ListenAndServe(":8082", mux)
}
```
