package telemetry

import (
	"context"
	"example/hello/config"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

var Tracer = otel.Tracer("default-tracer")

func handleErr(err error, message string) {
	if err != nil {
		GetLogger().Fatal().Err(err).Msg(message)
	}
}

func InitTraceProvider() func() {
	ctx := context.Background()

	config := config.ReadConfig()

	conn, err := grpc.DialContext(ctx, config.Tracing.Target, grpc.WithInsecure(), grpc.WithBlock())
	handleErr(err, "failed to create gRPC connection to collector")

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	handleErr(err, "failed to create trace exporter")

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String("0.0.1"),
	)

	// Aggregate spans before export
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(tp.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}

func TraceMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("gin-tracer")
}
