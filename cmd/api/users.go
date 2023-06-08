package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sulavmhrzn/chitrapost/internal/data"
)

type JWTCutsomClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

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
		return echo.NewHTTPError(code, echo.Map{
			"error": message,
		})
	}

	u := &data.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := data.ValidateUser(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	user, err := app.models.UserModel.GetUser(u)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
				"error": "invalid credentials",
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
				"error": "internal server error",
			})
		}
	}

	ok, err := data.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"error": "invalid credentials",
		})
	}
	claims := &JWTCutsomClaims{
		user.ID,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(app.cfg.JWT_SECRET))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})

}
