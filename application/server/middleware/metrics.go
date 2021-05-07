package middleware

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var FooCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "foo_total",
	Help: "Number of foo successfully processed.",
})

var Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "hits",
}, []string{"status", "path"})

func OkResponse(ctx echo.Context) {
	Hits.WithLabelValues("200", ctx.Request().URL.String()).Inc()
	FooCount.Add(1)
}

func ErrResponse(ctx echo.Context, code int) {
	Hits.WithLabelValues(strconv.Itoa(code), ctx.Request().URL.String()).Inc()
}
