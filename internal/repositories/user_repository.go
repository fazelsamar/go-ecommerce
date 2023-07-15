package repositories

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/pkg/database"
)

type UserRepository interface {
	Create(*entity.User) (*entity.User, error)
	GetByUsername(string) (*entity.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return userRepository{}
}

func (ur userRepository) Create(input *entity.User) (*entity.User, error) {
	db := database.GetDatabaseConnection()
	result := db.Create(&input)
	if result.Error != nil {
		return nil, result.Error
	}
	return input, nil
}

func (ur userRepository) GetByUsername(username string) (*entity.User, error) {
	db := database.GetDatabaseConnection()
	user := new(entity.User)
	result := db.Where(&entity.User{Username: username}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
