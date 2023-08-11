package global

import "os"

const (
	TracerKey = "otel-go-contrib-tracer"
)

var (
	// OTEL_ENDPOINT = "grpc://127.0.0.1:55680"
	OTEL_ENDPOINT = os.Getenv("OTEL_ENDPOINT")
)

func init() {
	if OTEL_ENDPOINT == "" {
		OTEL_ENDPOINT = "http://127.0.0.1:55681"
	}
}

var (
	AppName       = os.Getenv("AppName")
	App2_Endpoint = os.Getenv("App2_Endpoint")
	App4_Endpoint = os.Getenv("App4_Endpoint")
)

func init() {
	if AppName == "" {
		AppName = "OTEL_App1"
	}

	if App2_Endpoint == "" {
		App2_Endpoint = "http://127.0.0.1:3000/api/v1/app2/balance"
	}

	if App4_Endpoint == "" {
		App4_Endpoint = "http://127.0.0.1:3000/api/v1/app4/cellphone"
	}
}
