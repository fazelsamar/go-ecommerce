package models

import (
	"os"
	"time"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func Register(c *fiber.Ctx) error {
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
	user.Password = string(hash)
	user.IsAdmin = false
	db := database.DBConn
	result := db.Create(&user)
	if result.Error != nil {
		return c.Status(400).SendString("Cant create user")
	}

	// respond
	return c.Status(201).SendString("Created")
}

func Login(c *fiber.Ctx) error {
	// get username/password from req body
	reqUser := new(User)
	if err := c.BodyParser(reqUser); err != nil {
		return c.Status(400).SendString("Invalid body")
	}

	// get the user by username
	dbUser := new(User)
	db := database.DBConn
	db.Where(&User{Username: reqUser.Username}).First(&dbUser)
	if dbUser.ID <= 0 {
		return c.Status(400).SendString("Invalid username or password")
	}

	// compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password))
	if err != nil {
		return c.Status(400).SendString("User not found")
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  dbUser.ID,
		"isAdmin": dbUser.IsAdmin,
		"exp":     time.Now().Add(time.Hour * 24 * 30 * 365).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(400).SendString("Cant create token: " + err.Error() + tokenString)
	}

	// respond
	return c.JSON(fiber.Map{
		"token":   tokenString,
		"isAdmin": dbUser.IsAdmin,
	})
}
