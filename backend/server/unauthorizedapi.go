// Package server it is used for unauthorized apis
package server

import (
	"eventify/datatype"
	"eventify/db"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Secret Secret key imported from .env
var Secret = "Eventify"

// Unauthorized contains the api routes which don't need cookie
func Unauthorized(app *fiber.App) {
	app.Post("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		var body datatype.User
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if body.Password == "" || body.Email == "" || body.Name == "" || (body.Role != "organizer" && body.Role != "user") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password,Email,Name,Role might be missing"})
		}
		if body.ID == "" {
			body.ID = uuid.NewString()
		}

		password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		body.Password = string(password)
		claims := jwt.MapClaims{
			"id":  body.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		result, err := db.AddUsers(body)
		fmt.Println(result)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(Secret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Cookie(&fiber.Cookie{
			Name:     "sessional_id",
			Value:    tokenString,
			HTTPOnly: true,
			Secure:   false,
			SameSite: fiber.CookieSameSiteStrictMode,
			Expires:  time.Now().Add(time.Hour * 72),
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		password := c.Query("password")
		if password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password is missing"})
		}
		tokenString := c.Cookies("sessional_id")
		token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
			return []byte(Secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		claims := token.Claims.(jwt.MapClaims)
		id, ok := claims["id"].(string)
		if !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		body, err := db.FindUser(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		}
		if err = bcrypt.CompareHashAndPassword([]byte(body.Password), []byte(password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(body)
	})
}
