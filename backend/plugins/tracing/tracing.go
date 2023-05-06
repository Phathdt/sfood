package tracing

import (
	"context"
	"flag"

	sctx "github.com/viettranx/service-context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type tracingClient struct {
	id          string
	serviceName string
	jaegerHost  string
	logger      sctx.Logger
	tp          *tracesdk.TracerProvider
}

func NewTracingClient(id string, serviceName string) *tracingClient {
	tracer = otel.Tracer(serviceName)

	return &tracingClient{id: id, serviceName: serviceName}
}

func (t *tracingClient) ID() string {
	return t.id
}

func (t *tracingClient) InitFlags() {
	flag.StringVar(&t.jaegerHost, "jaeger_host", "http://localhost:14268/api/traces", "jaeger host")
}

func (t *tracingClient) Activate(sc sctx.ServiceContext) error {
	t.logger = sc.Logger(t.id)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(t.jaegerHost)))
	if err != nil {
		return err
	}

	t.tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(t.serviceName),
		)),
	)

	otel.SetTracerProvider(t.tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return nil
}

func (t *tracingClient) Stop() error {
	return t.tp.Shutdown(context.Background())
}
