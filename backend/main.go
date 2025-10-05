// package main contains the server configurations
package main

import (
	"eventify/db"
	"eventify/server"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
func main() {
	db.Connect()
	defer db.Disconnect()
	app := fiber.New(fiber.Config{
		AppName: "Eventify",
	})
	app.Use(limiter.New(limiter.Config{
		Max: 20,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,PUT,DELETE",
		AllowOrigins: "*",
	}))
	server.Unauthorized(app)
	server.Authorized(app)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})
	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
