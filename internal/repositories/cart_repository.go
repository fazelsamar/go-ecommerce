package repositories

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/google/uuid"
)

type CartRepository interface {
	Create(*entity.Cart) (*entity.Cart, error)
	GetCartItemsByCartId(uuid.UUID) ([]entity.CartItem, error)
}

type cartRepository struct {
}

func NewCartRepository() CartRepository {
	return cartRepository{}
}

func (cr cartRepository) Create(input *entity.Cart) (*entity.Cart, error) {
	db := database.GetDatabaseConnection()
	tx := db.Save(&input)
	return input, tx.Error
}

func (cr cartRepository) GetCartItemsByCartId(cart_id uuid.UUID) ([]entity.CartItem, error) {
	db := database.GetDatabaseConnection()
	var items []entity.CartItem
	tx := db.Preload("Product").Where("cart_id = ?", cart_id).Find(&items)
	return items, tx.Error
}
