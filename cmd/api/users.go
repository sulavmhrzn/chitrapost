package main

import (
	"errors"
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
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
				"error": err.Error(),
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(http.StatusCreated, user)
}

func (app *application) LoginUserHandler(c echo.Context) error {
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

	user, err := app.models.UserModel.GetUser(u)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
				"error": "invalid credentials",
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
				"error": "internal server error",
			})
		}
	}

	ok, err := data.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
	}
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]any{
			"error": "invalid credentials",
		})
	}
	return c.JSON(http.StatusOK, user)

}
