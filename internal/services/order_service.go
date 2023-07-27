package services

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/repositories"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(entity.Order) (*entity.Order, error)
	CreateOrderItems([]entity.OrderItem) error
	GetOrderItemsByOrderId(id uuid.UUID) ([]entity.OrderItem, error)
	GetOrderByUserId(user_id uint) ([]entity.Order, error)
}

type orderService struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return orderService{
		orderRepository: orderRepository,
	}
}

func (os orderService) CreateOrder(input entity.Order) (*entity.Order, error) {
	return os.orderRepository.Create(input)
}

func (os orderService) CreateOrderItems(input []entity.OrderItem) error {
	return os.orderRepository.CreateOrderItems(input)
}

func (os orderService) GetOrderItemsByOrderId(id uuid.UUID) ([]entity.OrderItem, error) {
	return os.orderRepository.GetOrderItemsByOrderId(id)
}

func (os orderService) GetOrderByUserId(user_id uint) ([]entity.Order, error) {
	return os.orderRepository.GetOrderByUserId(user_id)
}
