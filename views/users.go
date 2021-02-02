package views

import (
	"ConchaAPI/database"
	"ConchaAPI/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func CreateUserEndpoint(c *fiber.Ctx) error {
	user := &models.User{}
	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.SendString("Invalid user object")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return c.SendString("Invalid password")
	}

	user.Password = string(pass)
	createdUser := database.DBConn.Create(user)

	if createdUser.Error != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString(createdUser.Error.Error())
	}
	var createdJson []byte
	createdJson, err = json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(createdJson)
	return c.Send(createdJson)
}

func GetUsersEndpoint(c *fiber.Ctx) error {
	users, err := models.GetUsers()
	if err != nil {
		return c.SendString(err.Error())
	}
	var usersJson []byte
	usersJson, err = json.Marshal(users)
	if err != nil {
		fmt.Println(err)
	}
	return c.Send(usersJson)
}
