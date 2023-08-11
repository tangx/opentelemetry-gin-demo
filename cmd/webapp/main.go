package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tangx/opentelemetry-gin-demo/global"
	"github.com/tangx/opentelemetry-gin-demo/pkg/apis"
	"github.com/tangx/opentelemetry-gin-demo/pkg/apis/app1"
	"github.com/tangx/opentelemetry-gin-demo/pkg/apis/app2"
	"github.com/tangx/opentelemetry-gin-demo/pkg/apis/app4"
	"github.com/tangx/opentelemetry-gin-demo/pkg/apis/user"
	"github.com/tangx/opentelemetry-gin-demo/pkg/middlewares/otel"
)

func main() {
	r := gin.Default()

	r.Use(
		otel.Register(global.AppName, global.OTEL_ENDPOINT),
		otel.ReponseTraceID(),
	)

	v1 := r.Group("/api/v1")
	v1.Handle(http.MethodGet, "/ping", apis.Ping)

	v1.Handle(http.MethodGet, "/user/info", user.Info)

	// app tree
	{
		v1.Handle(http.MethodGet, "/app1/info", app1.UserInfoHandler)
		v1.Handle(http.MethodGet, "/app2/balance", app2.UserBalanceHandler)
		v1.Handle(http.MethodGet, "/app4/cellphone", app4.UserCellphoneHandler)

	}
	execute(r)
}

func execute(r *gin.Engine) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
