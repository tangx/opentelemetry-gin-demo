package utils

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

func MapCarrier(ctx context.Context) map[string]string {
	// 6. 向后传递 Header: traceparent
	pp := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)

	carrier := propagation.MapCarrier{}
	pp.Inject(ctx, carrier)

	return carrier
}
