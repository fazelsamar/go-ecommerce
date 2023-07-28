package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler interface {
	CreateOrder(*fiber.Ctx) error
	GetTheOrders(*fiber.Ctx) error
	GetOrderSerializer(entity.Order) ResponseOrder
	GetOrdersSerializerByUserId(entity.User) ResponseOrders
}

type orderHandler struct {
	orderService   services.OrderService
	cartService    services.CartService
	productService services.ProductService
}

func NewOrderHandler(
	orderService services.OrderService,
	cartService services.CartService,
	productService services.ProductService,
) OrderHandler {
	return orderHandler{
		orderService:   orderService,
		cartService:    cartService,
		productService: productService,
	}
}

type ResponseOrderItem struct {
	Product   entity.Product `json:"product"`
	Quantity  uint           `json:"quantity"`
	UnitPrice uint           `json:"unit_price"`
}

type ResponseOrder struct {
	ID        uuid.UUID           `json:"id"`
	CreatedAt time.Time           `json:"created_at"`
	Items     []ResponseOrderItem `json:"items"`
}

type ResponseOrders struct {
	Orders []ResponseOrder `json:"orders"`
}

func (oh orderHandler) GetOrderSerializer(order entity.Order) ResponseOrder {
	items, _ := oh.orderService.GetOrderItemsByOrderId(order.ID)
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

func (oh orderHandler) GetOrdersSerializerByUserId(user entity.User) ResponseOrders {
	orders, _ := oh.orderService.GetOrderByUserId(user.ID)
	orders_ser := ResponseOrders{Orders: make([]ResponseOrder, len(orders))}
	for order_index, order := range orders {
		items, _ := oh.orderService.GetOrderItemsByOrderId(order.ID)
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
		orders_ser.Orders[order_index] = order_ser
	}
	return orders_ser
}

func (oh orderHandler) CreateOrder(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(entity.User)

	// Check the cart
	id := c.Params("cart_id")
	cart, err := oh.cartService.GetCartById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"massage": "Cart Not Found!"})
		} else {
			return c.Status(500).JSON(fiber.Map{"massage": "Something went wrong!"})
		}
	}

	// Load the cartitems and check the count of items
	items, _ := oh.cartService.GetCartItemsByCartId(cart.ID)
	if len(items) <= 0 {
		return c.Status(400).JSON(fiber.Map{"massage": "Cart has no item"})
	}

	// Create the order and check the orderitems inventory
	var order entity.Order
	order.ID = uuid.New()
	order.UserID = user.ID
	orderItems := make([]entity.OrderItem, len(items))
	for index, item := range items {
		if item.Quantity > item.Product.Inventory {
			msg := fmt.Sprintf("Cartitem with productId=%d has not enough inventory", item.Product.ID)
			return c.Status(400).JSON(fiber.Map{"massage": msg})
		}
		orderItems[index].OrderID = order.ID
		orderItems[index].ProductID = item.Product.ID
		orderItems[index].Quantity = item.Quantity
		orderItems[index].UnitPrice = item.Product.Price
	}

	// Save the order
	order_res, err := oh.orderService.CreateOrder(order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Failed to create order!"})
	}

	// save the orderitems
	err = oh.orderService.CreateOrderItems(orderItems)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"massage": "Failed to create order items!"})
	}

	// Delete the cart
	oh.cartService.DeleteCartById(cart.ID)

	// TODO: Decrease the quantity of each product

	return c.JSON(oh.GetOrderSerializer(*order_res))
}

func (oh orderHandler) GetTheOrders(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(entity.User)

	return c.JSON(oh.GetOrdersSerializerByUserId(user))
}
