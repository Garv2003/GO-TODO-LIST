package routes

import (
	"github.com/Garv2003/TODOLIST/handler"
	"github.com/labstack/echo/v4"
)

func Setup(app *echo.Echo) {
	app.GET("/register", handler.GetRegister)
	app.POST("/register", handler.PostRegister)

	app.GET("/login", handler.GetLogin)
	app.POST("/login", handler.PostLogin)

	app.POST("/logout", handler.Logout)

	app.GET("/", handler.Home)

	app.DELETE("/delete/:id", handler.DeleteTodo)

	app.POST("/add", handler.AddTodo)
	app.POST("/toggle/:id", handler.IsComplete)
	app.POST("/edit/:id", handler.EditToDo)
}
