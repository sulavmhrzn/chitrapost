package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sulavmhrzn/chitrapost/internal/data"
)

type config struct {
	dsn string
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
		dsn: os.Getenv("DSN"),
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

	e.Use(middleware.Logger())

	api := e.Group("/api")
	api.GET("/healthcheck", HealtCheckHandler)

	usersGroup := api.Group("/users")
	usersGroup.POST("/register", app.RegisterUserHandler)

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
