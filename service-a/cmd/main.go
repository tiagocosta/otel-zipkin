package main

import (
	"fmt"

	"github.com/tiagocosta/otel-zipkin-service-a/configs"
	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web"
	"github.com/tiagocosta/otel-zipkin-service-a/internal/infra/web/webserver"
)

func main() {
	cfg := configs.LoadConfig[configs.Conf](".")

	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webZipCodeHandler := web.NewWebZipCodeHandler()
	webserver.AddHandler("/weather", webZipCodeHandler.ProcessZipCode)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	webserver.Start()
}
