package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	CreatedAt time.Time `json:"created_at"`
	CartID    uuid.UUID `json:"-" gorm:"primaryKey;index"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:CartID;references:ID"`
	ProductID uint      `json:"-" gorm:"primaryKey;index"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint      `json:"quantity"`
}

type Cart struct {
	ID        uuid.UUID      `json:"id" gorm:"index;primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ResponseCartItem struct {
	CreatedAt time.Time `json:"created_at"`
	Product   Product   `json:"product"`
	Quantity  uint      `json:"quantity"`
}

type ResponseCart struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Items     []ResponseCartItem
}

func NewCart(c *fiber.Ctx) error {
	db := database.GetDatabaseConnection()
	cart := new(Cart)
	cart.ID = uuid.New()
	db.Create(&cart)
	return c.JSON(cart)
}

func GetCartSerializer(cart Cart, db *gorm.DB) ResponseCart {
	var items []CartItem
	db.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items)
	// fmt.Println(items)
	response_items := make([]ResponseCartItem, len(items))
	for index, item := range items {
		response_items[index].CreatedAt = item.CreatedAt
		response_items[index].Product = item.Product
		response_items[index].Quantity = item.Quantity
	}
	cart_ser := ResponseCart{
		ID:        cart.ID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		DeletedAt: cart.DeletedAt,
		Items:     response_items,
	}
	return cart_ser
}

func AddItem(c *fiber.Ctx) error {
	// Check the cart
	id := c.Params("id")
	db := database.GetDatabaseConnection()
	var cart Cart
	if cart_result := db.First(&cart, "id = ?", id); cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Cart Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	// Get the product and quantity from the request body
	type RequestBody struct {
		ProductId uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
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
		return c.Status(400).SendString("Not enough inventory for product_id = " + fmt.Sprint(product.ID))
	}

	// Check the existing product in cart
	var cartItem CartItem
	var saveResult *gorm.DB
	db.Where("cart_id = ?", cart.ID).Where("product_id = ?", requestBody.ProductId).Find(&cartItem)
	if cartItem.Quantity == 0 {
		// Add product to cartItem
		cartItem = CartItem{
			CartID:    cart.ID,
			ProductID: requestBody.ProductId,
			Quantity:  requestBody.Quantity,
		}
		saveResult = db.Create(&cartItem)
	} else {
		// Check the quantity of product in cart and inventory of product
		if requestBody.Quantity+cartItem.Quantity > product.Inventory {
			return c.Status(400).SendString("Not enough inventory for product_id = " + fmt.Sprint(product.ID))
		}

		// Update quantity of product in cartItem
		cartItem.Quantity += requestBody.Quantity
		saveResult = db.Save(&cartItem)
	}

	if saveResult.Error != nil {
		return c.Status(500).SendString("Failed to add item to cart!")
	}

	return c.JSON(GetCartSerializer(cart, db))
}

func GetCart(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.GetDatabaseConnection()
	var cart Cart
	cart_result := db.First(&cart, "id = ?", id)
	if cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Cart Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}
	return c.JSON(GetCartSerializer(cart, db))
}

func DeleteCart(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.GetDatabaseConnection()
	var cart Cart
	cart_result := db.First(&cart, "id = ?", id)
	if cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Cart Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	db.Delete(&cart)
	return c.Status(204).SendString("")
}
