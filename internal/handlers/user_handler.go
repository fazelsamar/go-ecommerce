package handlers

import (
	"os"
	"time"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler interface {
	Register(*fiber.Ctx) error
	Login(*fiber.Ctx) error
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return userHandler{
		userService: userService,
	}
}

func (uh userHandler) Register(c *fiber.Ctx) error {
	// get username/password from req body
	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"massage": "Invalid body"})
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Cant hash password: " + err.Error()})
	}

	// create user
	user.Password = string(hash)
	user.IsAdmin = false
	_, err = uh.userService.CreateUser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"massage": "User already exists"})
	}

	return c.Status(201).JSON(fiber.Map{"massage": "Created"})
}

func generateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":  user.ID,
		"isAdmin": user.IsAdmin,
		"exp":     time.Now().Add(time.Hour * 24 * 30 * 365).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func (uh userHandler) Login(c *fiber.Ctx) error {
	// get username/password from req body
	reqUser := new(entity.User)
	if err := c.BodyParser(reqUser); err != nil {
		return c.Status(400).JSON(fiber.Map{"massage": "Invalid body"})
	}

	// get the user by username
	user, err := uh.userService.GetUserByUsername(reqUser.Username)
	if err != nil || user.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"massage": "User not found"})
	}

	// compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"massage": "User not found"})
	}

	// create token
	token, err := generateToken(user)
	if err != nil || len(token) == 0 {
		return c.Status(500).JSON(fiber.Map{"massage": "Cant create token: " + err.Error()})
	}

	// respond
	return c.Status(200).JSON(fiber.Map{
		"token":   token,
		"isAdmin": user.IsAdmin,
	})
}
