package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sulavmhrzn/chitrapost/internal/data"
)

func (app *application) UploadFileHandler(c echo.Context) error {
	url := make(chan string)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCutsomClaims)
	userID := claims.ID

	u, err := app.models.UserModel.GetUserFromID(userID)
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

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	go func() {
		resp, err := app.cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{Folder: "chitrapost"})
		if err != nil {
			log.Println(err)
			echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
				"error": "internal server error",
			})
		}
		url <- resp.URL
	}()
	input := &data.Chitra{
		URL:    <-url,
		UserID: u.ID,
	}

	chitra, err := app.models.ChitraModel.Insert(input)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"chitra": chitra,
	})
}
