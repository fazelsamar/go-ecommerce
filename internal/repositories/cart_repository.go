package repositories

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
	"github.com/google/uuid"
)

type CartRepository interface {
	Create(*entity.Cart) (*entity.Cart, error)
	GetCartItemsByCartId(uuid.UUID) ([]entity.CartItem, error)
	GetById(string) (entity.Cart, error)
	DeleteById(uuid.UUID) error
	GetCartItemByCartIdAndProductId(uuid.UUID, uint) (entity.CartItem, error)
	SaveCartItem(entity.CartItem) (entity.CartItem, error)
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

func (cr cartRepository) GetById(id string) (entity.Cart, error) {
	db := database.GetDatabaseConnection()
	var cart entity.Cart
	tx := db.Where("id = ?", id).First(&cart)
	return cart, tx.Error
}

func (cr cartRepository) DeleteById(id uuid.UUID) error {
	db := database.GetDatabaseConnection()
	var cart entity.Cart
	tx := db.Where("id = ?", id).Delete(&cart)
	return tx.Error
}

func (cr cartRepository) GetCartItemByCartIdAndProductId(cart_id uuid.UUID, product_id uint) (entity.CartItem, error) {
	db := database.GetDatabaseConnection()
	var items entity.CartItem
	tx := db.Preload("Product").Where("cart_id = ?", cart_id).Where("product_id = ?", product_id).Find(&items)
	return items, tx.Error
}

func (cr cartRepository) SaveCartItem(input entity.CartItem) (entity.CartItem, error) {
	db := database.GetDatabaseConnection()
	tx := db.Save(&input)
	return input, tx.Error
}
