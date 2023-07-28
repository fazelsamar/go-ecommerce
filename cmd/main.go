package main

import (
	"os"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/routes"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/fazelsamar/go-ecommerce/pkg/environment"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	environment.LoadEnvVariables()
	database.InitDB()
	db := database.GetDatabaseConnection()
	db.AutoMigrate(&entity.Product{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Cart{})
	db.AutoMigrate(&entity.CartItem{})
	db.AutoMigrate(&entity.Order{})
	db.AutoMigrate(&entity.OrderItem{})

	r := routes.NewRouting()

	// Use the CORS middleware with default options
	r.App.Use(cors.New())

	// Serve static files from the "mediafiles" directory
	r.App.Static("/mediafiles", "./mediafiles")

	r.App.Listen(":" + os.Getenv("EXPOSE_PORT"))
}
