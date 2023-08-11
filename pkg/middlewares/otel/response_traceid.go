package otel

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

func ReponseTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanctx, span := utils.Span(c, "Response Propagation")
		if span == nil {
			c.Next()
			return
		}
		defer span.End()

		// 4. 应答客户端时， 在 Header 中默认添加 TraceID
		traceid := span.SpanContext().TraceID().String()
		c.Header("TraceID", traceid)

		// 6. 向后传递 Header: traceparent
		pp := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
		)

		carrier := propagation.MapCarrier{}
		pp.Inject(spanctx, carrier)

		for k, v := range carrier {
			c.Header(k, v)
		}
	}
}
