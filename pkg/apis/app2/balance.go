package app2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

func UserBalanceHandler(c *gin.Context) {
	b := getBalance(c)
	c.String(http.StatusOK, fmt.Sprint(b))
}

func getBalance(ctx context.Context) string {

	ctx, span := utils.Span(ctx, "Get Balance")
	if span != nil {
		defer span.End()
	}

	b, _ := getBalanceFromDB(ctx)
	return b
}

func getBalanceFromDB(ctx context.Context) (balance string, err error) {
	_, span := utils.Span(ctx, "Get Balance From DB")
	if span != nil {
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "")
			}
			span.End()
		}()
	}

	time.Sleep(234 * time.Millisecond)

	if time.Now().Unix()%3 == 0 {
		return "unkown balance", fmt.Errorf("unkown balance")
	}
	return "1000", nil
}
