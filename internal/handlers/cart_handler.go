package handlers

import (
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
)

type CartHandler interface {
	NewCart(*fiber.Ctx) error
}

type cartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) CartHandler {
	return cartHandler{
		cartService: cartService,
	}
}

func (ch cartHandler) NewCart(c *fiber.Ctx) error {
	cart, err := ch.cartService.CreateCart()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Cant create cart: " + err.Error()})
	}
	return c.Status(200).JSON(cart)
}
