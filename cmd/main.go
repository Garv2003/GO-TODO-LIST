package main

import (
	"github.com/Garv2003/TODOLIST/db"
	"github.com/Garv2003/TODOLIST/routes"
	"github.com/labstack/echo/v4"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	db.Connect()
	app := echo.New()

	routes.Setup(app)

	app.Logger.Fatal(app.Start(":4000"))
}
