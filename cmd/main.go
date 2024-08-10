package main

import (
	"github.com/Garv2003/TODOLIST/db"
	"github.com/Garv2003/TODOLIST/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	db.Connect()

	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.BodyLimit("2M"))
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	app.Use(middleware.Secure())

	routes.Setup(app)

	app.Logger.Fatal(app.Start(":4000"))
}
