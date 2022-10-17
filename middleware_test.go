package echoprometheus

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestMetric(t *testing.T) {
	//arrange
	e := echo.New()
	e.Use(MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	request := httptest.NewRequest(echo.GET, "/metrics", nil)
	recorder := httptest.NewRecorder()

	//action
	e.ServeHTTP(recorder, request)

	//assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	assert.True(t, strings.Contains(recorder.Body.String(), "echo"))
	assert.True(t, strings.Contains(recorder.Body.String(), "http"))
}

func TestMetricWithCustomConfig(t *testing.T) {
	//arrange
	e := echo.New()
	var configMetrics = NewConfig()
	configMetrics.Namespace = "namespace"
	configMetrics.Subsystem = "test_subsystem"
	configMetrics.NormalizeHTTPStatus = false
	configMetrics.Buckets = []float64{
		1, // 1s
		2, // 2s
	}
	e.Use(MetricsMiddlewareWithConfig(configMetrics))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	request := httptest.NewRequest(echo.GET, "/metrics", nil)
	recorder := httptest.NewRecorder()

	//action
	e.ServeHTTP(recorder, request)

	//assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())
	assert.True(t, strings.Contains(recorder.Body.String(), "namespace"))
	assert.True(t, strings.Contains(recorder.Body.String(), "test_subsystem"))
}
