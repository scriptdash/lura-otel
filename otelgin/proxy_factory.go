package otelgin

import (
	"context"

	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/proxy"
	"go.opentelemetry.io/otel/trace"
)

func NewProxyFactory(factory proxy.Factory) proxy.Factory {
	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := factory.New(cfg)
		if err != nil {
			return proxy.NoopProxy, err
		}
		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {
			span, _ := ctx.Value(currentSpanKey).(trace.Span)
			if span != nil {
				ctx = trace.ContextWithSpan(ctx, span)
			}
			return next(ctx, request)
		}, nil
	})
}
