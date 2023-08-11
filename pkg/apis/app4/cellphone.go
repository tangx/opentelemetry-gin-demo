package app4

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

func UserCellphoneHandler(c *gin.Context) {

	phone, err := cellphoneFromRedis(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, phone)
}

func cellphoneFromRedis(ctx context.Context) (string, error) {

	_, span := utils.Span(ctx, "Cellphone From Database")
	if span != nil {
		defer span.End()
	}

	time.Sleep(500 * time.Millisecond)

	return "139-0000-1111", nil
}
