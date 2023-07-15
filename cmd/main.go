package main

import (
	"github.com/fazelsamar/go-ecommerce/internal/routes"
	"github.com/fazelsamar/go-ecommerce/models"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/fazelsamar/go-ecommerce/pkg/environment"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	environment.LoadEnvVariables()
	database.InitDB()
	db := database.GetDatabaseConnection()
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartItem{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.OrderItem{})

	r := routes.NewRouting()

	// Use the CORS middleware with default options
	r.App.Use(cors.New())

	// Serve static files from the "mediafiles" directory
	r.App.Static("/mediafiles", "./mediafiles")

	r.App.Listen(":8000")
}
