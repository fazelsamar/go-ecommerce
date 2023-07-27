package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	// get the token of the header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"massage": "Must authenticate"})
	}
	tokenString = strings.Split(string(tokenString), " ")[1]

	// Decode/validate the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(401).JSON(fiber.Map{"massage": "Token expired"})
		}

		// find the user
		var user entity.User
		db := database.GetDatabaseConnection()
		db.First(&user, claims["userId"])

		if user.ID == 0 {
			return c.Status(401).JSON(fiber.Map{"massage": "User not found"})
		}

		// attach the token to the request
		c.Locals("user", user)

		//continue
		return c.Next()
	} else {
		return c.Status(401).JSON(fiber.Map{"massage": "Invalid token"})
	}

}
