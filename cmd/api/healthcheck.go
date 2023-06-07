package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealtCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"environment": "development",
	})
}
