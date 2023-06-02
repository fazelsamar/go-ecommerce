package models

import (
	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

func NewUser(c *fiber.Ctx) error {
	// get username/password from req body
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString("Invalid body")
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(400).SendString("Cant hash password")
	}

	// create user
	db := database.DBConn
	user.Password = string(hash)
	result := db.Create(&user)
	if result.Error != nil {
		return c.Status(400).SendString("Cant create user")
	}

	// respond
	return c.Status(201).SendString("Created")
}
