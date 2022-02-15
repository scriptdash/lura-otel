package otelgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/proxy"
	router "github.com/luraproject/lura/router/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type contextKey int

const (
	tracerKey      string = "lura-otel-tracer"
	currentSpanKey string = "lura-otel-current-span"
)

func NewHandlerFactory(hf router.HandlerFactory) router.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		return HandlerFunc(cfg, hf(cfg, p))
	}
}

func HandlerFunc(cfg *config.EndpointConfig, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := func(_ http.ResponseWriter, r *http.Request) {
			span := trace.SpanFromContext(r.Context())
			c.Set(currentSpanKey, span)
			next(c)
		}
		otelHandler := otelhttp.NewHandler(
			h,
			c.Request.URL.Path,
			otelhttp.WithPropagators(propagation.TraceContext{}),
		)
		otelHandler.ServeHTTP(c.Writer, c.Request)
	}
}
