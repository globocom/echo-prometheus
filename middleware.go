package api

import (
	"reflect"

	echo "github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

)

var (
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "echo",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Number of HTTP operations",
	}, []string{"status", "method", "handler"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "echo",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Spend time by processing a route",
		Buckets: []float64{
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
		},
	}, []string{"method", "handler"})
)


func normalizeHTTPStatus(status int) string {
	if status < 200 {
		return "1xx"
	} else if status < 300 {
		return "2xx"
	} else if status < 400 {
		return "3xx"
	} else if status < 500 {
		return "4xx"
	}
	return "5xx"
}


func isNotFoundHandler(handler echo.HandlerFunc) bool {
	return reflect.ValueOf(handler).Pointer() == reflect.ValueOf(echo.NotFoundHandler).Pointer()
}

func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		path := c.Path()

		// to avoid attack high cardinality of 404
		if isNotFoundHandler(c.Handler()) {
			path = "/not-found"
		}

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(req.Method, path))
		err := next(c)
		timer.ObserveDuration()

		if err != nil {
			c.Error(err)
		}

		status := normalizeHTTPStatus(c.Response().Status)
		httpRequests.WithLabelValues(status, req.Method, path).Inc()

		return err
	}
}
