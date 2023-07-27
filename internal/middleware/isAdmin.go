package middleware

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func IsAdmin(c *fiber.Ctx) error {
	// check the user
	user, ok := c.Locals("user").(entity.User)
	if !ok || user.ID == 0 {
		return c.Status(500).JSON(fiber.Map{"massage": "User not found"})
	}

	// check for isAdmin
	if !user.IsAdmin {
		return c.Status(500).JSON(fiber.Map{"massage": "You must be admin"})
	}

	//continue
	return c.Next()
}
