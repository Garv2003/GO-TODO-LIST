package handler

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

//func Register(c *fiber.Ctx) error {
//	var data map[string]string
//	if err := c.BodyParser(&data); err != nil {
//		return err
//	}
//
//	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
//
//	user := models.User{Name: data["name"], Email: data["email"], Password: password}
//
//	db.DB.Create(&user)
//
//	return c.JSON(data)
//}
//
//func Login(c *fiber.Ctx) error {
//	var data map[string]string
//	if err := c.BodyParser(&data); err != nil {
//		return err
//	}
//
//	var user models.User
//	db.DB.Where("email = ?", data["email"]).First(&user)
//
//	if user.Id == 0 {
//		c.Status(fiber.StatusNotFound)
//		return c.JSON(fiber.Map{"message": "User not found"})
//	}
//
//	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
//		c.Status(fiber.StatusUnauthorized)
//		return c.JSON(fiber.Map{"message": "Incorrect password"})
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":  strconv.Itoa(int(user.Id)),
//		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
//	})
//
//	hmacSampleSecret := []byte("dndknskdnlsklkdslk")
//
//	tokenString, err := token.SignedString(hmacSampleSecret)
//
//	if err != nil {
//		c.Status(fiber.StatusInternalServerError)
//		return c.JSON(fiber.Map{"message": err.Error()})
//	}
//
//	cookie := fiber.Cookie{Name: "token", Value: tokenString, Expires: time.Now().Add(time.Hour * 24), HTTPOnly: true}
//
//	c.Cookie(&cookie)
//
//	return c.JSON(fiber.Map{"token": tokenString})
//}
//
//func User(c *fiber.Ctx) error {
//
//	cookie := c.Cookies("token")
//
//	fmt.Print(cookie)
//
//	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
//		// Don't forget to validate the alg is what you expect:
//		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//		//	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		//}
//		hmacSampleSecret := []byte("dndknskdnlsklkdslk")
//		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
//		return hmacSampleSecret, nil
//	})
//
//	if err != nil {
//		fmt.Println(err)
//		c.Status(fiber.StatusUnauthorized)
//		return c.JSON(fiber.Map{"message": "Incorrect tokendsds"})
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//
//	if !ok {
//		fmt.Println(ok)
//		c.Status(fiber.StatusUnauthorized)
//		return c.JSON(fiber.Map{"message": "Incorrect token"})
//	}
//
//	var user models.User
//
//	db.DB.Where("id = ?", claims["id"]).First(&user)
//
//	return c.JSON(fiber.Map{"user": user})
//}

func GetRegister(c *fiber.Ctx) error {

}

func PostRegister(c *fiber.Ctx) error {

}

func GetLogin(c *fiber.Ctx) error {}

func PostLogin(c *fiber.Ctx) error {}

func Home(c *fiber.Ctx) error {}

func AddTodo(c *fiber.Ctx) error {}

func DeleteTodo(c *fiber.Ctx) error {}

func IsComplete(c *fiber.Ctx) error {}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"message": "Success"})
}
