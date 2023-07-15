package repositories

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
)

type ProductRepository interface {
	Create(*entity.Product) (*entity.Product, error)
	GetAll() ([]entity.Product, error)
	GetById(uint) (entity.Product, error)
	Delete(*entity.Product) error
}

type productRepository struct {
}

func NewProductRepository() ProductRepository {
	return productRepository{}
}

func (pr productRepository) Create(input *entity.Product) (*entity.Product, error) {
	db := database.GetDatabaseConnection()
	tx := db.Save(&input)
	return input, tx.Error
}

func (ur productRepository) GetAll() ([]entity.Product, error) {
	db := database.GetDatabaseConnection()
	var products []entity.Product
	tx := db.Order("id DESC").Find(&products)
	return products, tx.Error
}

func (ur productRepository) GetById(id uint) (entity.Product, error) {
	db := database.GetDatabaseConnection()
	var product entity.Product
	tx := db.First(&product, id)
	return product, tx.Error
}

func (ur productRepository) Delete(input *entity.Product) error {
	db := database.GetDatabaseConnection()
	tx := db.Delete(&input)
	return tx.Error
}
