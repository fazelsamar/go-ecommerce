package handlers

import (
	"errors"
	"fmt"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartHandler interface {
	NewCart(*fiber.Ctx) error
	AddItem(*fiber.Ctx) error
}

type cartHandler struct {
	cartService    services.CartService
	productService services.ProductService
}

func NewCartHandler(
	cartService services.CartService,
	productService services.ProductService,
) CartHandler {
	return cartHandler{
		cartService:    cartService,
		productService: productService,
	}
}

func (ch cartHandler) NewCart(c *fiber.Ctx) error {
	cart, err := ch.cartService.CreateCart()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Cant create cart: " + err.Error()})
	}
	return c.Status(200).JSON(cart)
}

func (ch cartHandler) AddItem(c *fiber.Ctx) error {
	// Check the cart
	id := c.Params("id")
	cart, err := ch.cartService.GetCartById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"massage": "Cart Not Found!"})
		} else {
			return c.Status(500).JSON(fiber.Map{"massage": "Something went wrong!"})
		}
	}

	// Get the product and quantity from the request body
	type RequestBody struct {
		ProductId uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
	}
	var requestBody RequestBody
	err = c.BodyParser(&requestBody)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"massage": "Invalid request body!"})
	}

	// Check the product
	product, err := ch.productService.GetProductById(requestBody.ProductId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"massage": "Product not Found!"})
		} else {
			return c.Status(500).JSON(fiber.Map{"massage": "Something went wrong!"})
		}
	}

	// Check the product quantity
	if requestBody.Quantity > product.Inventory {
		return c.Status(400).JSON(fiber.Map{"massage": "Not enough inventory for product_id = " + fmt.Sprint(product.ID)})
	}

	// Check the existing product in cart
	cartItem, _ := ch.cartService.GetCartItemByCartIdAndProductId(cart.ID, requestBody.ProductId)
	if cartItem.Quantity == 0 {
		// Add product to cartItem
		cartItem = entity.CartItem{
			CartID:    cart.ID,
			ProductID: requestBody.ProductId,
			Quantity:  requestBody.Quantity,
		}
		cartItem, err = ch.cartService.SaveCartItem(cartItem)
	} else {
		// Check the quantity of product in cart and inventory of product
		if requestBody.Quantity+cartItem.Quantity > product.Inventory {
			return c.Status(400).SendString("Not enough inventory for product_id = " + fmt.Sprint(product.ID))
		}

		// Update quantity of product in cartItem
		cartItem.Quantity += requestBody.Quantity
		cartItem, err = ch.cartService.SaveCartItem(cartItem)
	}

	if err != nil {
		return c.Status(500).SendString("Failed to add item to cart!")
	}

	return c.JSON(ch.cartService.GetCartSerializer(cart))
}
