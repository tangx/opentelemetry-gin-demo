package app1

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"github.com/tangx/opentelemetry-gin-demo/global"
	"github.com/tangx/opentelemetry-gin-demo/pkg/httpclient"
	"github.com/tangx/opentelemetry-gin-demo/pkg/utils"
)

type UserInfo struct {
	Name      string
	Cellphone string
	Balance   int
}

type Response struct {
	Data  any
	Error string
}

func UserInfoHandler(c *gin.Context) {

	name := c.GetHeader("UserName")
	uinfo, err := userInfo(c, name)
	if err != nil {
		resp := &Response{
			Error: err.Error(),
		}

		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := &Response{
		Data: uinfo,
	}

	c.JSON(http.StatusOK, resp)
}

func userInfo(ctx context.Context, name string) (*UserInfo, error) {

	uinfo := &UserInfo{
		Name: name,
	}

	b, err := balanceFromApp2(ctx)
	if err != nil {
		return nil, err
	}
	uinfo.Balance = b

	p, err := cellphoneFromApp4(ctx)
	if err != nil {
		return nil, err
	}
	uinfo.Cellphone = p

	return uinfo, nil
}

func balanceFromApp2(ctx context.Context) (result int, err error) {

	opt := trace.WithSpanKind(trace.SpanKindClient)
	ctx, span := utils.Span(ctx, "Get balance from next app", opt)
	if span != nil {

		defer func() {

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			span.End()
		}()

	}

	// time cost
	time.Sleep(100 * time.Millisecond)

	data, err := httpclient.GET(ctx, global.App2_Endpoint)
	if err != nil {
		return 0, err
	}

	b, err := strconv.Atoi(data)
	if err != nil {
		return 0, err
	}

	return b, nil

}

func cellphoneFromApp4(ctx context.Context) (phone string, err error) {

	ctx, span := utils.Span(ctx, "Cellphone from Next App")
	if span != nil {
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

			span.End()
		}()
	}

	// lock := mutex.Lock()
	lock := &sync.Mutex{}

	span.AddEvent("Get CellPhone: Start")
	lock.Lock()

	span.AddEvent("Get CellPhone: Request")
	phone, err = httpclient.GET(ctx, global.App4_Endpoint)
	if err != nil {
		return "", err
	}

	span.AddEvent("Get Cellphone: Verify")
	lock.Unlock()

	if IsValidPhone(phone) {
		return phone, nil
	}

	err = fmt.Errorf("Error: invalid cellphone, %s", phone)
	return "", err
}

var (
	patt = `\d{3}-?\d{4}-?\d{4}`
	reg  = regexp.MustCompile(patt)
)

func IsValidPhone(phone string) bool {
	return reg.Match([]byte(phone))
}
