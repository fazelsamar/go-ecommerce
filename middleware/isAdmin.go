package middleware

import (
	"fmt"

	"github.com/fazelsamar/go-ecommerce/models"
	"github.com/gofiber/fiber/v2"
)

func IsAdmin(c *fiber.Ctx) error {
	// check the user
	user, ok := c.Locals("user").(models.User)
	fmt.Println(user.ID)
	if !ok || user.ID == 0 {
		return c.Status(500).SendString("")
	}

	// check for isAdmin
	if !user.IsAdmin {
		return c.Status(400).SendString("You must be admin")
	}

	//continue
	return c.Next()
}
