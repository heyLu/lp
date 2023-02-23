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

## Monitor(ing) without `calls_total` and `latency_bucket`

This is just an idea, but maybe worth playing around with for systems with lots of existing metrics that already have most of the required info?

- metrics only use [range queries](https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries) ([source](https://github.com/jaegertracing/jaeger/blob/main/plugin/metrics/prometheus/metricsstore/reader.go))
- would need to "find" or map existing metrics
  - different metric name
  - different label names
  - missing labels (`span_kind`, `status_code`)
    - generate based on defaults ("SPAN_KIND_SERVER") and `http_status_code`
  - looks like we need an "on the fly" Prometheus proxy that rewrites queries as needed
- there is <https://pkg.go.dev/github.com/prometheus/prometheus/promql/parser>, which may be exactly what is needed?
- query example: `curl 'http://localhost:9090/api/v1/query_range' --get --data-urlencode 'query=sum(rate(calls_total{service_name =~ "example-service", status_code = "STATUS_CODE_ERROR"}[5m])) by (service_name) / sum(rate(calls_total{service_name =~ "example-service"}[5m])) by (service_name)' --data-urlencode 'start=1677186016' --data-urlencode 'end=1677186500' --data-urlencode 'step=10s'`
  - would need rewriting of `data.result.metric` which contains labels (i.e. `service_name` because of the aggregation?)
  - could decode using `json.Decoder#Token`
    - but need to output json ourselves
    - need to figure out how to find out that we are at `metric`
    - or just parse in-memory :)

## Resources

- https://www.jaegertracing.io/docs/1.42/spm/
- https://pkg.go.dev/github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanmetricsprocessor#section-readme
