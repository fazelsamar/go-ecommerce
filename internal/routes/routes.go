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
	userservice := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userservice)

	productRepo := repositories.NewProductRepository()
	productservice := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productservice)

	cartRepo := repositories.NewCartRepository()
	cartservice := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartservice, productservice)

	orderRepo := repositories.NewOrderRepository()
	orderservice := services.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderservice, cartservice, productservice)

	// Auth routes
	r.App.Post("/register", userHandler.Register)
	r.App.Post("/login", userHandler.Login)

	// // Product routes
	r.App.Get("/products", productHandler.GetProductList)
	r.App.Get("/products/:id", productHandler.GetProduct)
	r.App.Post("/products", middleware.RequireAuth, middleware.IsAdmin, productHandler.NewProduct)
	r.App.Delete("/products/:id", middleware.RequireAuth, middleware.IsAdmin, productHandler.DeleteProduct)

	// // Cart routes
	r.App.Get("/cart", cartHandler.NewCart)
	r.App.Get("/cart/:id", cartHandler.GetCart)
	r.App.Post("/cart/:id", cartHandler.AddItem)
	r.App.Delete("/cart/:id", cartHandler.DeleteCart)

	// // Order routes
	r.App.Get("/order/:cart_id", middleware.RequireAuth, orderHandler.CreateOrder)
	r.App.Get("/orders", middleware.RequireAuth, orderHandler.GetTheOrders)
}
