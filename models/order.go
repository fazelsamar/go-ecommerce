package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fazelsamar/go-ecommerce/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id" gorm:"index;primaryKey"`
	UserID    uint      `json:"-" gorm:"index"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

type OrderItem struct {
	OrderID   uuid.UUID `json:"-" gorm:"primaryKey;index"`
	Order     Order     `json:"order" gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint      `json:"-" gorm:"primaryKey;index"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  uint      `json:"quantity"`
	UnitPrice uint      `json:"unit_price"`
}

type ResponseOrderItem struct {
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
	UnitPrice uint    `json:"unit_price"`
}

type ResponseOrder struct {
	ID        uuid.UUID           `json:"id"`
	CreatedAt time.Time           `json:"created_at"`
	Items     []ResponseOrderItem `json:"items"`
}

type ResponseOrders struct {
	Orders []ResponseOrder `json:"orders"`
}

func GetOrderSerializer(order Order, db *gorm.DB) ResponseOrder {
	var items []OrderItem
	db.Preload("Product").Where("order_id = ?", order.ID).Find(&items)
	// fmt.Println(items)
	response_items := make([]ResponseOrderItem, len(items))
	for index, item := range items {
		response_items[index].Product = item.Product
		response_items[index].Quantity = item.Quantity
		response_items[index].UnitPrice = item.UnitPrice
	}
	order_ser := ResponseOrder{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		Items:     response_items,
	}
	return order_ser
}

func GetOrdersSerializerByUserId(user User, db *gorm.DB) ResponseOrders {
	var orders []Order
	db.Where("user_id = ?", user.ID).Find(&orders)
	orders_ser := ResponseOrders{Orders: make([]ResponseOrder, len(orders))}
	for order_index, order := range orders {
		var items []OrderItem
		db.Preload("Product").Where("order_id = ?", order.ID).Find(&items)
		response_items := make([]ResponseOrderItem, len(items))
		for index, item := range items {
			response_items[index].Product = item.Product
			response_items[index].Quantity = item.Quantity
			response_items[index].UnitPrice = item.UnitPrice
		}
		fmt.Println(response_items)
		order_ser := ResponseOrder{
			ID:        order.ID,
			CreatedAt: order.CreatedAt,
			Items:     response_items,
		}
		orders_ser.Orders[order_index] = order_ser
	}
	return orders_ser
}

func CreateOrder(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(User)

	// Check the cart
	cartId := c.Params("cart_id")
	db := database.GetDatabaseConnection()
	var cart Cart
	if cart_result := db.First(&cart, "id = ?", cartId); cart_result.Error != nil {
		if errors.Is(cart_result.Error, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("Cart Not Found!")
		} else {
			return c.Status(500).SendString("Something went wrong!")
		}
	}

	// Load the cartitems and check the count of items
	var items []CartItem
	db.Preload("Product").Where("cart_id = ?", cart.ID).Find(&items)
	if len(items) <= 0 {
		return c.Status(400).SendString("Cart has no item")
	}

	// Create the order and check the orderitems inventory
	var order Order
	order.ID = uuid.New()
	order.UserID = user.ID
	orderItems := make([]OrderItem, len(items))
	for index, item := range items {
		if item.Quantity > item.Product.Inventory {
			msg := fmt.Sprintf("Cartitem with productId=%d has not enough inventory", item.Product.ID)
			return c.Status(400).SendString(msg)
		}
		orderItems[index].OrderID = order.ID
		orderItems[index].ProductID = item.Product.ID
		orderItems[index].Quantity = item.Quantity
		orderItems[index].UnitPrice = item.Product.Price
	}

	// Save the order
	res := db.Create(&order)
	if res.Error != nil {
		return c.Status(500).SendString("Failed to create order!")
	}

	// save the orderitems
	res = db.Create(&orderItems)
	if res.Error != nil {
		return c.Status(500).SendString("Failed to create order items!")
	}

	// Delete the cart
	db.Delete(&cart)

	return c.JSON(GetOrderSerializer(order, db))
}

func GetTheOrders(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(User)
	db := database.GetDatabaseConnection()
	return c.JSON(GetOrdersSerializerByUserId(user, db))
}
