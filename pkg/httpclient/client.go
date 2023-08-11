package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

func GET(ctx context.Context, url string) (string, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	headers := utils.MapCarrier(ctx)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// client := http.DefaultClient
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
