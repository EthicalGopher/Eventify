// Package server it is used for unauthorized apis
package server

import (
	"eventify/db"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User type for storing user data
type User struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	Interests []string `bson:"interests,omitempty"`
	Role      string   `bson:"role"`
}

// Secret Secret key imported from .env
var Secret = "Eventify"

// Unauthorized contains the api routes which don't need cookie
func Unauthorized(app *fiber.App) {
	app.Post("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		var body db.User
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if body.ID == "" {
			body.ID = uuid.NewString()
		}

		password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		body.Password = string(password)
		claims := jwt.MapClaims{
			"id":  body.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		db.AddUsers(body)
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
			SameSite: "Strict",
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
	})
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}
