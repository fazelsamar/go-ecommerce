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

	// Auth routes
	app.Post("/register", models.Register)
	app.Post("/login", models.Login)

	// Cart routes
	app.Get("/cart", models.NewCart)
	app.Get("/cart/:id", models.GetCart)
	app.Post("/cart/:id", models.AddItem)
	app.Delete("/cart/:id", models.DeleteCart)

	// Order routes
	app.Get("/order/:cart_id", middleware.RequireAuth, models.CreateOrder)
	app.Get("/orders", middleware.RequireAuth, models.GetTheOrders)
}

func main() {
	initializers.LoadEnvVariables()
	database.InitDB()
	db := database.GetDatabaseConnection()
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartItem{})
	db.AutoMigrate(&models.Order{})
	db.AutoMigrate(&models.OrderItem{})

	app := fiber.New()

	// Use the CORS middleware with default options
	app.Use(cors.New())

	// Serve static files from the "mediafiles" directory
	app.Static("/mediafiles", "./mediafiles")

	setupRoutes(app)

	app.Listen(":8000")
}
