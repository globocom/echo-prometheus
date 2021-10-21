package main

import (
	"net/http"

	echoPrometheus "github.com/globocom/echo-prometheus/v2"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()

	ec := echoPrometheus.MetricsMiddleware()
	ec.
		WithBuckets([]float64{0.01, 0.1, 0.5, 1.0, 5.0, 10.0}).
		NormalizeHTTPStatus(true).
		WithNamespace("MyCustomNamespace").
		WithSubsystem("MyCustomSubsystem")

	e.Use(ec.MetricsMiddlewareFunc())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
