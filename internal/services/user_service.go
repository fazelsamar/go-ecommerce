package services

import (
	"github.com/fazelsamar/go-ecommerce/internal/entity"
	"github.com/fazelsamar/go-ecommerce/internal/repositories"
)

type UserService interface {
	CreateUser(*entity.User) (*entity.User, error)
	GetUserByUsername(string) (*entity.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return userService{
		userRepository: userRepository,
	}
}

func (us userService) CreateUser(input *entity.User) (*entity.User, error) {
	return us.userRepository.Create(input)
}

func (us userService) GetUserByUsername(username string) (*entity.User, error) {
	return us.userRepository.GetByUsername(username)
}
