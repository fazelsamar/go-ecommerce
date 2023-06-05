package main

import (
	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/fazelsamar/go-ecommerce/initializers"
	"github.com/fazelsamar/go-ecommerce/middleware"
	"github.com/fazelsamar/go-ecommerce/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func setupRoutes(app *fiber.App) {
	// Product routes
	app.Get("/products", models.GetProducts)
	app.Get("/products/:id", models.GetProduct)
	app.Post("/products", middleware.RequireAuth, middleware.IsAdmin, models.NewProducts)
	app.Delete("/products/:id", middleware.RequireAuth, middleware.IsAdmin, models.DeleteProducts)

	//Auth routes
	app.Post("/register", models.Register)
	app.Post("/login", models.Login)
}

func main() {
	initializers.LoadEnvVariables()
	database.InitDB()
	database.DBConn.AutoMigrate(&models.Product{})
	database.DBConn.AutoMigrate(&models.User{})

	app := fiber.New()

	// Use the CORS middleware with default options
	app.Use(cors.New())

	// Serve static files from the "mediafiles" directory
	app.Static("/mediafiles", "./mediafiles")

	setupRoutes(app)

	app.Listen(":8000")
}
