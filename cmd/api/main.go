package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sulavmhrzn/chitrapost/internal/data"
)

type config struct {
	dsn        string
	JWT_SECRET string
}
type application struct {
	models data.Models
	cfg    config
}

func main() {
	if err := loadENV(); err != nil {
		log.Fatal(err)
	}
	cfg := config{
		dsn:        os.Getenv("DSN"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}
	conn, err := OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())
	app := application{
		cfg:    cfg,
		models: *data.NewModels(conn),
	}

	e := echo.New()
	jwtconfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCutsomClaims)
		},
		SigningKey: []byte(app.cfg.JWT_SECRET),
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())

	api := e.Group("/api")
	api.GET("/healthcheck", HealtCheckHandler, echojwt.WithConfig(jwtconfig))

	usersGroup := api.Group("/users")
	usersGroup.POST("/register", app.RegisterUserHandler)
	usersGroup.POST("/login", app.LoginUserHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func OpenDB(cfg config) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, cfg.dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

func loadENV() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
