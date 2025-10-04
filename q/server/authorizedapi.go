// Package server is used for authorized api
package server

import "github.com/gofiber/fiber/v2"

// Authorized Pending
func Authorized(app *fiber.App) {
	app.Get("/events", Auth, func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

// Auth Pending
func Auth(c *fiber.Ctx) error {
	return c.Next()
}
