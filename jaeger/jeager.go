package jaeger

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/ExcitingFrog/go-core-common/log"
	"github.com/ExcitingFrog/go-core-common/provider"
)

const (
	JaegerTrace = "jaeger_trace"
)

var globalTracer trace.Tracer

type Jaeger struct {
	provider.IProvider

	Config *Config
	Tracer trace.Tracer
	tp     *tracesdk.TracerProvider
}

func NewJaeger(config *Config) *Jaeger {
	if config == nil {
		config = NewConfig()
	}
	return &Jaeger{
		Config: config,
	}
}

func (j *Jaeger) Run() error {
	j.Tracer = otel.Tracer(j.Config.ServiceName)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(j.Config.JaegerURI)))
	if err != nil {
		return err
	}
	j.tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(j.Config.ServiceName),
		)),
	)
	otel.SetTracerProvider(j.tp)
	globalTracer = j.Tracer
	// otel.SetTextMapPropagator(propagation.TraceContext{})
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)
	return nil
}

func (j *Jaeger) Close() error {
	return j.tp.Shutdown(context.Background())
}

func (j *Jaeger) Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return j.Tracer.Start(ctx, name)
}

func StartSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Span) {
	if value := ctx.Value(JaegerTrace); value != nil {
		if span, ok := value.(trace.Span); ok {
			return ctx, span
		}
	}
	return globalTracer.Start(ctx, operationName)
}

func StartSpanAndLogFromContext(ctx context.Context, operationName string) (context.Context, trace.Span, *zap.Logger) {
	if value := ctx.Value(JaegerTrace); value != nil {
		if span, ok := value.(trace.Span); ok {
			return ctx, span, log.Logger()
		}
	}
	ctx, span := globalTracer.Start(ctx, operationName)
	return ctx, span, log.Logger()
}
