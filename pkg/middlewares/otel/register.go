package otel

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register(appname string, endpoint string) gin.HandlerFunc {

	opts := []otelgin.Option{
		ProviderOption(appname, endpoint),
		PropagationExtractOption(),
	}

	return otelgin.Middleware(appname, opts...)
}

func ProviderOption(appname string, endpoint string) otelgin.Option {
	// 1. 注册 Provider
	provider, err := initProvider(appname, endpoint)
	if err != nil {
		panic(err)
	}

	return otelgin.WithTracerProvider(provider)
}
