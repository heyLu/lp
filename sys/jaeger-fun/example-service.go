package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const ServiceName = "example-service"

var config struct {
	Addr string
}

// metrics for jaeger monitor view

var (
	CallsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "calls_total",
	}, []string{"http_method", "http_status_code", "operation", "service_name", "span_kind", "status_code"})
	LatencyBucket = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "latency",
		Buckets: []float64{2, 5, 10, 50, 100, 250, 500, 1000},
	}, []string{"http_method", "http_status_code", "operation", "service_name", "span_kind", "status_code"})
)

// tracing

var tracer trace.Tracer

func newExporter(ctx context.Context) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint())
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	flag.StringVar(&config.Addr, "addr", "localhost:3000", "Address to listen on")

	prometheus.MustRegister(
		CallsTotal,
		LatencyBucket,
	)

	exp, err := newExporter(context.Background())
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)

	defer func() { _ = tp.Shutdown(context.Background()) }()

	otel.SetTracerProvider(tp)

	tracer = tp.Tracer("ExampleService")

	http.HandleFunc("/", recordStats(func(w http.ResponseWriter, req *http.Request) {
		_, span := tracer.Start(req.Context(), req.URL.Path)
		defer span.End()

		dur := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(dur)

		if rand.Float64() > 0.9 {
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Fprintf(w, "%s", dur)
	}))

	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{Registry: prometheus.DefaultRegisterer}))

	log.Printf("Listening on http://%s", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}

func recordStats(handlerFn func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		rw := &recordResponseWriter{ResponseWriter: w, statusCode: 200}
		handlerFn(rw, req)

		go func(dur time.Duration, method string, path string, statusCode int) {
			status := "STATUS_CODE_UNSET"
			if statusCode >= 400 {
				status = "STATUS_CODE_ERROR"
			}
			CallsTotal.WithLabelValues(method, fmt.Sprintf("%d", statusCode), path, ServiceName, "SPAN_KIND_SERVER", status).Inc()
			LatencyBucket.WithLabelValues(method, fmt.Sprintf("%d", statusCode), path, ServiceName, "SPAN_KIND_SERVER", status).Observe(float64(dur.Milliseconds()))
		}(time.Since(start), req.Method, req.URL.Path, rw.statusCode)
	}
}

type recordResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func (rrw *recordResponseWriter) WriteHeader(statusCode int) {
	rrw.statusCode = statusCode
	rrw.ResponseWriter.WriteHeader(statusCode)
}
