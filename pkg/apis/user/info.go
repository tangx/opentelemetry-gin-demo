package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"github.com/tangx/opentelemetry-gin-demo/pkg/httpclient"
	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

var (
	USER_INFO_HOST = os.Getenv("USER_INFO_HOST")
)

// Info 获取用户信息
// https://zhuanlan.zhihu.com/p/608282493
func Info(c *gin.Context) {

	username := c.GetHeader("UserName")
	if username == "" {
		username = "jane"
	}

	name := fmt.Sprintf("RequestURI: %s", c.Request.RequestURI)
	spanctx, span := utils.Span(c, name)
	defer span.End()

	data, err := info(spanctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))

		return
	}

	c.JSON(http.StatusOK, data)
}

func info(ctx context.Context, name string) (*UserInfo, error) {

	// 注入 attr 属性
	ctx = utils.SpanContextWithAttr(ctx, map[string]string{"user.name": name})

	// 设置为 consumer kind
	opt := trace.WithSpanKind(trace.SpanKindConsumer)

	spanctx, span := utils.Span(ctx, "user info integration", opt)
	if span != nil {
		defer span.End()
	}

	userinfo := &UserInfo{
		Name: name,
	}

	b, err := balance(spanctx, name)
	if err != nil {
		return nil, err
	}
	userinfo.Balance = b

	c, err := cellphone(spanctx, name)
	if err != nil {
		return nil, err
	}
	userinfo.Cellphone = c

	return userinfo, nil
}

// balance get user balance
func balance(ctx context.Context, name string) (int, error) {
	ctx = utils.SpanContextWithAttr(ctx, map[string]string{"user.kind": "func.balance"})

	_, span := utils.Span(ctx, "user balance")
	if span != nil {
		defer span.End()
	}

	switch name {
	case "guanyu":
		return 100, nil

	case "zhangfei":
		return 200, nil
	}

	return 0, errors.New("unknown user")
}

func cellphone(ctx context.Context, name string) (string, error) {
	ctx = utils.SpanContextWithAttr(ctx, map[string]string{"user.kind": "func.cellphone"})

	ctx, span := utils.Span(ctx, "user cellphone")
	if span != nil {
		defer span.End()
	}

	switch name {
	case "guanyu":
		return "131-1111-2222", nil
		// case "zhangfei":
		// 	return "132-2222-3333", nil
	}

	err := errors.New("unknown user or cellphone not found")

	// 提交错误日志
	span.RecordError(err)

	// 设置状态
	span.SetStatus(codes.Error, "unsupport user")

	attrs := semconv.HTTPAttributesFromHTTPStatusCode(500)
	span.SetAttributes(attrs...)

	// 设置属性
	// span.SetAttributes(attribute.KeyValue{
	// 	Key:   "user.kind",
	// 	Value: attribute.StringValue("user.cellphone"),
	// })

	if os.Getenv("PORT") != "9099" {
		httpclient.GET(ctx, "http://127.0.0.1:9099/api/v1/user/info")
	}

	return "", err
}

type UserInfo struct {
	Name      string
	Balance   int
	Cellphone string
}
