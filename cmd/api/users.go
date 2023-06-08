package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sulavmhrzn/chitrapost/internal/data"
)

func (app *application) RegisterUserHandler(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		code := err.(*echo.HTTPError).Code
		message := err.(*echo.HTTPError).Message
		return echo.NewHTTPError(code, map[string]any{
			"error": message,
		})
	}
	u := &data.User{
		Email:    input.Email,
		Password: input.Password,
	}
	if err := data.ValidateUser(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}
	hash, err := data.CreateHashPassword(u.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	u.Password = hash
	user, err := app.models.UserModel.Insert(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, user)
}
