# otel-zipkin

This is a project that uses OpenTelemetry and Zipkin to create tracing data between two services (A and B). Both services have a tracer that is initialized with them. Basically, the service A is the entry point. It receives a request with a zipcode, does some validations on it and sends a new request to service B, witch returns a response to service A with the weather on that zipcode.

As soon as service A receives a request it creates a span to monitor the time taken to process that request. When calling service B endpoint, service A propagates the trace context using its registered propagator. Service B, now, has access to the same context and can create spans from it. This allows a complete tracing between services A and B.

![otel-zipkin](https://github.com/tiagocosta/otel-zipkin/assets/807693/544e95bf-4623-4c48-8b84-c685bc35bd36)

## Dependencies
Go 1.22.2

OpenTelemetry

Zipkin
