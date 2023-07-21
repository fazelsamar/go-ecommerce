package services

import (
	"time"

	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartService interface {
	CreateCart() (*entity.Cart, error)
}

type cartService struct {
	cartRepository repositories.CartRepository
}

func NewCartService(cartRepository repositories.CartRepository) CartService {
	return cartService{
		cartRepository: cartRepository,
	}
}

type ResponseCartItem struct {
	CreatedAt time.Time      `json:"created_at"`
	Product   entity.Product `json:"product"`
	Quantity  uint           `json:"quantity"`
}

type ResponseCart struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Items     []ResponseCartItem
}

func (cs cartService) CreateCart() (*entity.Cart, error) {
	cart := new(entity.Cart)
	cart.ID = uuid.New()
	return cs.cartRepository.Create(cart)
}

func (cs cartService) GetCartSerializer(cart entity.Cart) interface{} {
	items, _ := cs.cartRepository.GetCartItemsByCartId(cart.ID)

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
