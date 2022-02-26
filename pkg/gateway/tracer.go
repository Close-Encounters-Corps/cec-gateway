package gateway

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerConfig struct {
	ServiceName string
	ServiceVer  string
	Jaeger      string
	Environment string
	Disabled    bool
}

type Tracer struct {
	provider trace.TracerProvider
}

func (g *Gateway) SetupTracing(cfg *TracerConfig) error {
	if cfg.Disabled {
		provider := trace.NewNoopTracerProvider()
		g.Tracer = &Tracer{provider: provider}
		return nil
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Jaeger)))
	if err != nil {
		return err
	}
	prv := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVer),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		)),
	)
	otel.SetTracerProvider(prv)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	g.Tracer = &Tracer{provider: prv}
	return nil
}

func (p Tracer) Close(ctx context.Context) error {
	if prv, ok := p.provider.(*sdktrace.TracerProvider); ok {
		return prv.Shutdown(ctx)
	}
	return nil
}
