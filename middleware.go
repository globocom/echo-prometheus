package echoprometheus

import (
	"reflect"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Config responsible to configure middleware
type Config struct {
	Namespace           string
	Buckets             []float64
	Subsystem           string
	NormalizeHTTPStatus bool
}

const (
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
	notFoundPath         = "/not-found"
)

type EchoPrometheus struct {
	config Config
}

// DefaultConfig has the default instrumentation config
var DefaultConfig = Config{
	Namespace: "echo",
	Subsystem: "http",
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
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	},
	NormalizeHTTPStatus: true,
}

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

// WithNamespace change the namespace
func (ep *EchoPrometheus) WithNamespace(namespace string) *EchoPrometheus {
	ep.config.Namespace = namespace
	return ep
}

// WithBuckets change the buckets
func (ep *EchoPrometheus) WithBuckets(buckets []float64) *EchoPrometheus {
	ep.config.Buckets = buckets
	return ep
}

// WithSubsystem change the subsystem
func (ep *EchoPrometheus) WithSubsystem(subsystem string) *EchoPrometheus {
	ep.config.Subsystem = subsystem
	return ep
}

// WithNormalizeHTTPStatus change the normalizeHTTPStatus flag
func (ep *EchoPrometheus) WithNormalizeHTTPStatus(normalizeHTTPStatus bool) *EchoPrometheus {
	ep.config.NormalizeHTTPStatus = normalizeHTTPStatus
	return ep
}

// MetricsMiddlewareWithConfig returns an echo middleware for instrumentation.
func (ep EchoPrometheus) MetricsMiddlewareFunc() echo.MiddlewareFunc {
	httpRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ep.config.Namespace,
		Subsystem: ep.config.Subsystem,
		Name:      httpRequestsCount,
		Help:      "Number of HTTP operations",
	}, []string{"status", "method", "handler"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ep.config.Namespace,
		Subsystem: ep.config.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "Spend time by processing a route",
		Buckets:   ep.config.Buckets,
	}, []string{"method", "handler"})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := c.Path()

			// to avoid attack high cardinality of 404
			if isNotFoundHandler(c.Handler()) {
				path = notFoundPath
			}

			timer := prometheus.NewTimer(httpDuration.WithLabelValues(req.Method, path))
			err := next(c)
			timer.ObserveDuration()

			if err != nil {
				c.Error(err)
			}

			status := ""
			if ep.config.NormalizeHTTPStatus {
				status = normalizeHTTPStatus(c.Response().Status)
			} else {
				status = strconv.Itoa(c.Response().Status)
			}

			httpRequests.WithLabelValues(status, req.Method, path).Inc()

			return err
		}
	}
}

// MetricsMiddleware returns an echo middleware with default config for instrumentation.
func MetricsMiddleware() *EchoPrometheus {
	return &EchoPrometheus{
		config: DefaultConfig,
	}
}
