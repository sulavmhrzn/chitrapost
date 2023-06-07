package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	api := e.Group("/api")
	api.GET("/healthcheck", HealtCheckHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
