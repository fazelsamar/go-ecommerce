package main

import (
	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/fazelsamar/go-ecommerce/initializers"
	"github.com/fazelsamar/go-ecommerce/middleware"
	"github.com/fazelsamar/go-ecommerce/models"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	// Product routes
	app.Get("/products", models.GetProducts)
	app.Get("/product/:id", models.GetProduct)
	app.Post("/product", middleware.RequireAuth, models.NewProducts)
	app.Delete("/product/:id", middleware.RequireAuth, models.DeleteProducts)

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

	setupRoutes(app)

	app.Listen(":3000")
}
