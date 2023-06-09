package models

import (
	"errors"
	"time"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	// gorm.Model
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CartItems []CartItem     `json:"cart_items"`
}

type CartItem struct {
	gorm.Model
	Cart      Cart      `json:"cart"`
	CartID    uuid.UUID `gorm:"type:uuid" gorm:"index" json:"cart_id"`
	Product   Product   `json:"product"`
	ProductID uint      `json:"product_id"`
	Quantity  uint      `json:"quantity"`
}

func NewCart(c *fiber.Ctx) error {
	db := database.DBConn
	cart := new(Cart)
	cart.ID = uuid.New()
	db.Create(&cart)
	return c.JSON(cart)
}

func Item(c *fiber.Ctx) error {
	// Check the cart
	id := c.Params("id")
	db := database.DBConn
	var cart Cart
	cart_result := db.First(&cart, "id = ?", id)
	if cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	// Get the product and quantity from the request body
	type RequestBody struct {
		ProductId uint
		Quantity  uint
	}
	var requestBody RequestBody
	err := c.BodyParser(&requestBody)
	if err != nil {
		return c.Status(400).SendString("Invalid request body!")
	}

	// Check the product
	var product Product
	product_result := db.First(&product, "id = ?", requestBody.ProductId)
	if product_result.Error != nil {
		if errors.Is(product_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Product not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	// Check the product quantity
	if requestBody.Quantity > product.Inventory {
		return c.Status(400).SendString("Not enough inventory!")
	}

	// Add product to cartItems
	cartItem := CartItem{
		CartID:    cart.ID,
		ProductID: requestBody.ProductId,
		Quantity:  requestBody.Quantity,
	}

	// Add the CartItem to Cart's CartItems slice
	cart.CartItems = append(cart.CartItems, cartItem)

	// Save the changes to the database
	saveResult := db.Save(&cart)
	if saveResult.Error != nil {
		return c.Status(500).SendString("Failed to add item to cart!")
	}

	return c.JSON(cart)
}

func GetCart(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var cart Cart
	cart_result := db.First(&cart, "id = ?", id)
	if cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	var cart_item []CartItem
	cart_item_result := db.Find(&cart_item, "cart_id = ?", cart.ID)
	if cart_item_result.Error != nil {
		if errors.Is(cart_item_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	return c.JSON(cart_item_result)
}
