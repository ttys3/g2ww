package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Blackhole(c echo.Context) error {
	// since we have BodyDump middleware, here we just do nothing
	return c.JSON(http.StatusOK, "ok")
}