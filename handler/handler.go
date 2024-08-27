package handler

import (
	"net/http"
	"time"

	"github.com/Garv2003/TODOLIST/db"
	"github.com/Garv2003/TODOLIST/models"
	"github.com/Garv2003/TODOLIST/view/components"
	"github.com/Garv2003/TODOLIST/view/pages"
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}
	return ctx.HTML(statusCode, buf.String())
}

func GetRegister(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Register())
}

func PostRegister(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if name == "" || email == "" || password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	userID := uuid.New().String()
	user := models.User{Id: userID, Name: name, Email: email, Password: hashPassword}
	if err := db.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	return c.Redirect(http.StatusSeeOther, "/login")
}

func GetLogin(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Login())
}

func PostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	db.DB.Where("email = ?", email).First(&user)
	if user.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Incorrect password"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create token"})

	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func Home(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	var User models.User
	db.DB.Where("id = ?", claims["id"]).First(&User)
	var TodoList []models.Todo
	db.DB.Where("user_id = ?", User.Id).Find(&TodoList)
	return Render(c, http.StatusOK, pages.Home(TodoList, "", User))
}

func AddTodo(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	var User models.User
	db.DB.Where("id = ?", claims["id"]).First(&User)
	var data map[string]string
	if err := c.Bind(&data); err != nil {
		return err
	}
	ToDoId := uuid.New().String()
	todo := models.Todo{Content: data["content"], IsCompleted: false, UserId: User.Id, Id: ToDoId}
	if err := db.DB.Create(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create todo"})
	}
	return Render(c, http.StatusOK, components.Todo(todo))
}

func DeleteTodo(c echo.Context) error {

	cookie, err := c.Cookie("token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	id := c.Param("id")
	var todo models.Todo
	db.DB.Where("id = ?", id).First(&todo)
	if todo.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Todo not found"})
	}
	if err := db.DB.Delete(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete todo"})
	}

	var count int64
	db.DB.Model(&models.Todo{}).Where("user_id = ?", claims["id"]).Count(&count)
	if count == 0 {
		return c.HTML(http.StatusOK, "<p class='text-center'>No Todos</p>")
	}
	return c.HTML(http.StatusOK, "")
}

func IsComplete(c echo.Context) error {
	id := c.Param("id")
	var todo models.Todo
	db.DB.Where("id = ?", id).First(&todo)
	if todo.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Todo not found"})
	}
	todo.IsCompleted = !todo.IsCompleted
	if err := db.DB.Save(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update todo"})
	}
	return Render(c, http.StatusOK, components.Todo(todo))
}

func EditToDo(c echo.Context) error {
	id := c.Param("id")
	var todo models.Todo
	db.DB.Where("id = ?", id).First(&todo)
	if todo.Id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Todo not found"})
	}
	Content := c.FormValue("content")
	todo.Content = Content
	if err := db.DB.Save(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update todo"})
	}
	return Render(c, http.StatusOK, components.Todo(todo))
}

func Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()

	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}
