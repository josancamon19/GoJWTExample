package auth

import (
	"ConchaAPI/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct{}

func (h *Handler) Auth(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := models.GetUserByEmail(email)
	if err != nil {
		c.Status(http.StatusNotFound)
		return c.SendString(err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.Status(http.StatusBadRequest)
		return c.SendString("Invalid login credentials. Please try again")
	}
	tokens, err := generateTokenPair(user)
	return c.JSON(tokens)
}

// This is the api to refresh tokens
// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	err := json.Unmarshal(c.Body(), &tokenReq)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return c.SendString("RefreshToken not sent")
	}

	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		userID := int(claims["sub"].(float64))
		var user *models.User
		user, err = models.GetUserByID(userID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.SendString("User not found")
		}

		newTokenPair, err := generateTokenPair(user)
		if err != nil {
			return err
		}
		return c.JSON(newTokenPair)
	}
	return err
}

// https://echo.labstack.com/cookbook/jwt
func (h *Handler) private(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name + "!")
}
