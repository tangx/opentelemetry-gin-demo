package otel

import (
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/propagation"
)

// PropagationExtractOption 从上游获取 traceparent, tracestate
func PropagationExtractOption() otelgin.Option {
	tc := propagation.TraceContext{}
	return otelgin.WithPropagators(tc)
}
