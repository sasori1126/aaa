package otel

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func InitOtel() (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := getResource(ctx)
	if err != nil {
		return nil, err
	}
	exporter, err := traceExporter(ctx)
	if err != nil {
		return nil, err
	}
	bsp := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return provider.Shutdown, nil
}

func traceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	grpcConn, err := gRpcConn(ctx)
	if err != nil {
		return nil, err
	}
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(grpcConn))
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func RecordSpanError(sn trace2.Span, err error) {
	if !sn.IsRecording() {
		sn.RecordError(err)
		sn.SetStatus(codes.Error, err.Error())
	}
}

func gRpcConn(ctx context.Context) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:4317",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Axisforestry"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "prod"),
		),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
