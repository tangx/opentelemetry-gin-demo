// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example using OTLP exporters + collector + third-party backens. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package otel

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(appname string, endpoint string) (trace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(appname),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if len(endpoint) == 0 {
		endpoint = "127.0.0.1:55680"
	}

	// var err error
	exporter, err := traceExporter(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	return tracerProvider, nil
}

func traceExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	ur, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var exporter *otlptrace.Exporter

	switch strings.ToLower(ur.Scheme) {
	case "http", "https":
		exporter, err = httpExporter(ctx, ur.Host)
		if err != nil {
			return nil, err
		}

	case "grpc":
		fallthrough
	default:
		exporter, err = grpcExpoter(ctx, ur.Host)
		if err != nil {
			return nil, err
		}
	}

	return exporter, nil
}

// 创建 OTEL 的 GRPC 连接器
func grpcExpoter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {
	// addr := strings.TrimLeft(endpoint, "grpc://")

	conn, err := grpc.DialContext(ctx, endpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		// grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
		// otlptracegrpc.WithHeaders(
		// 	map[string]string{
		// 		"authorization": BearerAuthToken,
		// 		"Authorization": BearerAuthToken,
		// 	},
		// ),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}

func httpExporter(ctx context.Context, endpoint string) (*otlptrace.Exporter, error) {

	// endpoint = strings.TrimPrefix(endpoint, "https://")
	// endpoint = strings.TrimPrefix(endpoint, "http://")

	opts := []otlptracehttp.Option{
		otlptracehttp.WithTimeout(5 * time.Second),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
		// otlptracehttp.WithHeaders(
		// 	map[string]string{
		// 		"authorization": BearerAuthToken,
		// 		"Authorization": BearerAuthToken,
		// 	},
		// ),
	}

	trace, err := otlptracehttp.New(ctx, opts...)

	return trace, err
}
