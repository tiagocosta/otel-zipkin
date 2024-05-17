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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var logger = log.New(os.Stderr, "zipkin-service-a ", log.Ldate|log.Ltime|log.Llongfile)

func initTracer(url string) (func(context.Context) error, error) {
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
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := configs.LoadConfig[configs.Conf](".")

	url := "http://zipkin-all-in-one:9411/api/v2/spans"

	shutdown, err := initTracer(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("service-a-tracer")
	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webZipCodeHandler := web.NewWebZipCodeHandler(tracer)
	webserver.AddHandler("/weather", webZipCodeHandler.ProcessZipCode)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go webserver.Start()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}
}
