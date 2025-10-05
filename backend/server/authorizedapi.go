// Package server is used for authorized api
package server

import (
	"eventify/datatype"
	"eventify/db"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GetID finds the id using the Cookies
func GetID(c *fiber.Ctx) (string, error) {
	tokenString := c.Cookies("sessional_id")
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
		return []byte(Secret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	id, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid token")
	}
	return id, nil
}

// Auth Pending
func auth(c *fiber.Ctx) error {
	id, err := GetID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = db.FindUser(id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Next()
}

// Authorized checks if the user is in the database
func Authorized(app *fiber.App) {
	// Add Event
	app.Get("/event/add", auth, func(c *fiber.Ctx) error {
		var event datatype.Event
		if err := c.BodyParser(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		OrganizerID, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		user, err := db.FindUser(OrganizerID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if user.Role != "organizer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not an organizer"})
		}
		event.OrganizerID = OrganizerID
		result, err := db.AddEvents(event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": result})
	})
	// Shows All events
	app.Get("/event", auth, func(c *fiber.Ctx) error {
		id, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		events, err := db.FindEvents(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(events)

	})
	// Edit event details
	app.Put("/event/edit", auth, func(c *fiber.Ctx) error {
		var event datatype.Event
		if err := c.BodyParser(&event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		err := db.EditEvent(event.ID, event)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})
	// Edit user details
	app.Put("/user/edit", auth, func(c *fiber.Ctx) error {
		var user datatype.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		id, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		user.ID = id
		if err := db.EditUser(id, user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Add user to event
	app.Post("/event/:id/join", auth, func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		userID, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.AddUserToEvent(eventID, userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Remove user from event
	app.Post("/event/:id/leave", auth, func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		userID, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.RemoveUserFromEvent(eventID, userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Delete user
	app.Delete("/user/delete", auth, func(c *fiber.Ctx) error {
		id, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.DeleteUser(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Delete event (organizer only)
	app.Delete("/event/:id/delete", auth, func(c *fiber.Ctx) error {
		eventID := c.Params("id")
		userID, err := GetID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		user, err := db.FindUser(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if user.Role != "organizer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not an organizer"})
		}
		if err := db.DeleteEvent(eventID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})
}
