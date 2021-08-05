package main

import (
	"go.uber.org/zap"
	"os"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	initLog()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		zap.S().Infow("begin body dump", "req", reqBody, "rsp", resBody)
	}))

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/stat", GwStat)

	// https://work.weixin.qq.com/api/doc/90000/90136/91770
	// 消息发送频率限制
	// 每个机器人发送的消息不能超过20条/分钟
	e.Any("/:key", GwWorker)

	// Start server
	var ListenAddress string
	if os.Getenv("DOCKER") != "" {
		ListenAddress = "0.0.0.0:2408"
	} else {
		ListenAddress = "127.0.0.1:2408"
	}
	zap.S().Infow("http server starting", "listen_addr", ListenAddress)
	e.Logger.Fatal(e.Start(ListenAddress))

}
