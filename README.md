# Echo Prometheus
Middleware for echo v4 to instrument all handlers as metrics


## Example of usage

### With default config
```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoPrometheus "github.com/globocom/echo-prometheus"
)

func main() {
	e := echo.New()

	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```

### With custom config
```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoPrometheus "github.com/globocom/echo-prometheus"
)

func main() {
	e := echo.New()

	var configMetrics = echoPrometheus.NewConfig()
		configMetrics.Namespace = "namespace"
		configMetrics.Buckets = []float64{
			0.0005, // 0.5ms
			0.001,  // 1ms
			0.005,  // 5ms
			0.01,   // 10ms
			0.05,   // 50ms
			0.1,    // 100ms
			0.5,    // 500ms
			1,      // 1s
			2,      // 2s
	}

	e.Use(echoPrometheus.MetricsMiddlewareWithConfig(configMetrics))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
```


## Example output for metric route

```
# HELP echo_http_request_duration_seconds Spend time by processing a route
# TYPE echo_http_request_duration_seconds histogram
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.0005"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.001"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.002"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.005"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.01"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.02"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.05"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.1"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.2"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="0.5"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="1"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="2"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="5"} 7
echo_http_request_duration_seconds_bucket{handler="/",method="GET",le="+Inf"} 7
echo_http_request_duration_seconds_sum{handler="/",method="GET"} 7.645099999999999e-05
echo_http_request_duration_seconds_count{handler="/",method="GET"} 7
# HELP echo_http_requests_total Number of HTTP operations
# TYPE echo_http_requests_total counter
echo_http_requests_total{handler="/",method="GET",status="2xx"} 7
```

### View metrics via Grafana

We built a grafana dashboard for these metrics, lookup at [https://grafana.com/grafana/dashboards/10913](https://grafana.com/grafana/dashboards/10913).

