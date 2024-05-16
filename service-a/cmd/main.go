package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/tiagocosta/otel-zipkin-service-a/configs"
	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web"
	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web/webserver"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

func initTracer(ctx context.Context) (func(context.Context) error, error) {
	url := "http://zipkin-all-in-one:9411/api/v2/spans"

	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("service-a"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := configs.LoadConfig[configs.Conf](".")

	// url := "http://localhost:9411/api/v2/spans"

	shutdown, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("microservice-tracer")
	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webZipCodeHandler := web.NewWebZipCodeHandler(tracer)
	webserver.AddHandler("/weather", webZipCodeHandler.ProcessZipCode)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	webserver.Start()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}
}

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"github.com/tiagocosta/otel-zipkin-service-a/configs"
// 	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web"
// 	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web/webserver"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
// )

// func main() {
// 	sigCh := make(chan os.Signal, 1)
// 	signal.Notify(sigCh, os.Interrupt)

// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer cancel()

// 	res, err := resource.New(ctx,
// 		resource.WithAttributes(
// 			semconv.ServiceName("mysrv1"),
// 		),
// 	)
// 	if err != nil {
// 		fmt.Errorf("failed to create resource: %w", err)
// 	}
// 	ctx, cancel = context.WithTimeout(ctx, time.Second)
// 	defer cancel()
// 	conn, err := grpc.DialContext(ctx, "otel-collector:4317",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		fmt.Errorf("failed to create gRPC connection to collector: %w", err)
// 	}

// 	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
// 	if err != nil {
// 		fmt.Errorf("failed to create trace exporter: %w", err)
// 	}

// 	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithResource(res),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)
// 	otel.SetTracerProvider(tracerProvider)

// 	otel.SetTextMapPropagator(propagation.TraceContext{})
// 	cfg := configs.LoadConfig[configs.Conf](".")
// 	tracer := otel.Tracer("microservice-tracer")
// 	webserver := webserver.NewWebServer(cfg.WebServerPort)
// 	webZipCodeHandler := web.NewWebZipCodeHandler(tracer)
// 	webserver.AddHandler("/weather", webZipCodeHandler.ProcessZipCode)
// 	fmt.Println("Starting web server on port", cfg.WebServerPort)
// 	webserver.Start()
// 	// // Create a span
// 	// ctx = context.Background()
// 	// ctx, span := tracer.Start(ctx, "iniciando")
// 	// for i := 0; i < 1000; i++ {
// 	// 	fmt.Println(i)
// 	// }
// 	// defer span.End()

// 	// ctx, span2 := tracer.Start(ctx, "iniciando 2")
// 	// time.Sleep(2 * time.Second)
// 	// defer span2.End()

// 	// // Add Baggage
// 	// member, _ := baggage.NewMember("username", "user123")
// 	// bag, _ := baggage.New(member)
// 	// ctx = baggage.ContextWithBaggage(ctx, bag)

// 	fmt.Println("Microservice 1")
// 	select {
// 	case <-sigCh:
// 		log.Println("Shutting down gracefully, CTRL+C pressed...")
// 	case <-ctx.Done():
// 		log.Println("Shutting down due to other reason...")
// 	}
// }
