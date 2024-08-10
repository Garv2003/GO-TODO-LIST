package routes

import (
	"github.com/Garv2003/TODOLIST/view"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Setup(app *echo.Echo) {
	app.GET("/", HomeHandler)
	//app.Get("/register", handler.GetRegister)
	//app.Post("/register", handler.PostRegister)
	//
	//app.Get("/login", handler.GetLogin)
	//app.Post("/login", handler.PostLogin)
	//
	//app.Get("/home", handler.Home)
	//
	//app.Delete("/delete/:id", handler.DeleteTodo)
	//
	//app.Post("/add", handler.AddTodo)
	//app.Post("/complete", handler.IsComplete)
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, app.App())
}
