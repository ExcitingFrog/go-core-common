package utrace

import (
	"context"

	"github.com/ExcitingFrog/go-core-common/log"
	"github.com/ExcitingFrog/go-core-common/provider"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var globalTracer trace.Tracer

type UTrace struct {
	provider.IProvider
	Config *Config
}

func NewUTrace(config *Config) *UTrace {
	if config == nil {
		config = NewConfig()
	}
	return &UTrace{
		Config: config,
	}
}

func (t *UTrace) Init() error {
	return nil
}

func (t *UTrace) Close() error {
	return uptrace.Shutdown(context.Background())
}

func (t *UTrace) Run() error {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(t.Config.UptraceDSN),
		uptrace.WithServiceName(t.Config.ServiceName),
		uptrace.WithServiceVersion("1.0.0"),
	)

	globalTracer = otel.Tracer(t.Config.ServiceName)

	return nil
}

func StartSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Span) {
	return globalTracer.Start(ctx, operationName)
}

func StartSpanAndLogFromContext(ctx context.Context, operationName string) (context.Context, trace.Span, *zap.Logger) {
	ctx, span := globalTracer.Start(ctx, operationName)
	return ctx, span, log.Logger()
}

func GetGlobalTracer() trace.Tracer {
	return globalTracer
}
