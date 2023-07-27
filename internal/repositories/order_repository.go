package repositories

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(entity.Order) (*entity.Order, error)
	CreateOrderItems([]entity.OrderItem) error
	GetOrderItemsByOrderId(uuid.UUID) ([]entity.OrderItem, error)
	GetOrderByUserId(uint) ([]entity.Order, error)
}

type orderRepository struct {
}

func NewOrderRepository() OrderRepository {
	return orderRepository{}
}

func (or orderRepository) Create(input entity.Order) (*entity.Order, error) {
	db := database.GetDatabaseConnection()
	tx := db.Create(&input)
	return &input, tx.Error
}

func (or orderRepository) CreateOrderItems(input []entity.OrderItem) error {
	db := database.GetDatabaseConnection()
	tx := db.Create(&input)
	return tx.Error
}

func (or orderRepository) GetOrderItemsByOrderId(id uuid.UUID) ([]entity.OrderItem, error) {
	db := database.GetDatabaseConnection()
	var items []entity.OrderItem
	tx := db.Preload("Product").First(&items, id)
	return items, tx.Error
}

func (or orderRepository) GetOrderByUserId(user_id uint) ([]entity.Order, error) {
	db := database.GetDatabaseConnection()
	var orders []entity.Order
	tx := db.Where("user_id = ?", user_id).Find(&orders)
	return orders, tx.Error
}
