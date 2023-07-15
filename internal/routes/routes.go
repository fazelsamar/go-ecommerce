package routes

import (
	"github.com/fazelsamar/go-ecommerce/internal/handlers"
	"github.com/fazelsamar/go-ecommerce/internal/middleware"
	"github.com/fazelsamar/go-ecommerce/internal/repositories"
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Routing struct {
	App *fiber.App
}

func NewRouting() *Routing {
	r := &Routing{
		App: fiber.New(),
	}
	r.routes()
	return r
}

func (r *Routing) routes() {

	userRepo := repositories.NewUserRepository()
	userUsecase := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	productRepo := repositories.NewProductRepository()
	productUsecase := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productUsecase)

	// Auth routes
	r.App.Post("/register", userHandler.Register)
	r.App.Post("/login", userHandler.Login)

	// // Product routes
	r.App.Get("/products", productHandler.GetProductList)
	r.App.Get("/products/:id", productHandler.GetProduct)
	r.App.Post("/products", middleware.RequireAuth, middleware.IsAdmin, productHandler.NewProduct)
	r.App.Delete("/products/:id", middleware.RequireAuth, middleware.IsAdmin, productHandler.DeleteProduct)

	// // Cart routes
	// r.App.Get("/cart", models.NewCart)
	// r.App.Get("/cart/:id", models.GetCart)
	// r.App.Post("/cart/:id", models.AddItem)
	// r.App.Delete("/cart/:id", models.DeleteCart)

	// // Order routes
	// r.App.Get("/order/:cart_id", middleware.RequireAuth, models.CreateOrder)
	// r.App.Get("/orders", middleware.RequireAuth, models.GetTheOrders)
}