package main

import (
	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/fazelsamar/go-ecommerce/product"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/products", product.GetProducts)
	app.Get("/product/:id", product.GetProduct)
	app.Post("/product", product.NewProducts)
	app.Delete("/product/:id", product.DeleteProducts)
}

func main() {
	app := fiber.New()

	database.InitDB()
	database.DBConn.AutoMigrate(&product.Product{})

	setupRoutes(app)

	app.Listen(":3000")
}
