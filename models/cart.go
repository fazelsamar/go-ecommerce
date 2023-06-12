package models

import (
	"errors"
	"time"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uuid.UUID `json:"-"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:CartID;references:ID;index"`
	ProductID uint      `json:"-"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint      `json:"quantity"`
}

type Cart struct {
	// gorm.Model
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// Items []CartItem `json:"items" gorm:"foreignKey:ID"`
}

type CartSer struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Items     []CartItem     `json:"items"`
}

func NewCart(c *fiber.Ctx) error {
	db := database.DBConn
	cart := new(Cart)
	cart.ID = uuid.New()
	db.Create(&cart)
	return c.JSON(cart)
}

func GetCartSerializer(cart Cart, db *gorm.DB) CartSer {
	var items []CartItem
	db.Where("cart_id = ?", cart.ID).Find(&items)

	cart_ser := CartSer{
		ID:        cart.ID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		DeletedAt: cart.DeletedAt,
		Items:     items,
	}
	return cart_ser
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
		// ProductID: product.ID,
		Quantity: requestBody.Quantity,
	}

	// Add the CartItem to Cart's CartItems slice
	// cart.Items = append(cart.Items, cartItem)

	// Save the changes to the database
	saveResult := db.Save(&cartItem)
	if saveResult.Error != nil {
		return c.Status(500).SendString("Failed to add item to cart!")
	}
	saveResult = db.Save(&cart)
	if saveResult.Error != nil {
		return c.Status(500).SendString("Failed to add item to cart!")
	}

	return c.JSON(GetCartSerializer(cart, db))
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
			return c.Status(500).SendString("Something went wrong!" + cart_result.Error.Error())
		}
	}

	return c.JSON(GetCartSerializer(cart, db))
}
