package services

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/repositories"
)

type ProductService interface {
	CreateProduct(*entity.Product) (*entity.Product, error)
	GetAllProducts() ([]entity.Product, error)
	GetProductById(uint) (entity.Product, error)
	DeleteProduct(*entity.Product) error
}

type productService struct {
	productRepository repositories.ProductRepository
}

func NewProductService(productRepository repositories.ProductRepository) ProductService {
	return productService{
		productRepository: productRepository,
	}
}

func (us productService) GetAllProducts() ([]entity.Product, error) {
	return us.productRepository.GetAll()
}

func (us productService) GetProductById(id uint) (entity.Product, error) {
	return us.productRepository.GetById(id)
}

func (us productService) CreateProduct(input *entity.Product) (*entity.Product, error) {
	return us.productRepository.Create(input)
}

func (us productService) DeleteProduct(input *entity.Product) error {
	return us.productRepository.Delete(input)
}
