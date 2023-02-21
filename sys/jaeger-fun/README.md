# jaeger-fun -- getting a little monitoring going

Tiny demo project of a service doing [OpenTelemetry](https://opentelemetry.io) tracing via
Jaeger, with that fancy [Service Performance Monitoring](https://www.jaegertracing.io/docs/1.42/spm/)
hacked in using two custom Prometheus metrics.

## Run?

- `docker-compose up` for the dependencies (jaeger and prometheus)
- `go run example-service.go` for the tiny service
- do some requests at <http://localhost:3000>
- have a look at <http://localhost:16686/monitor> for the SPM stuff (select `example-service`)
  - <http://localhost:9090> is the prometheus server

How neat!

## Resources

- https://www.jaegertracing.io/docs/1.42/spm/
- https://pkg.go.dev/github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor#section-readme
