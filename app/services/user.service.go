package services

import (
	"errors"

	"note-api/app/dto"
	"note-api/app/models"
	"note-api/app/repositories"
)

type UserService interface {
	CreateUser(user *dto.UserRequest) (*models.UserModel, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) CreateUser(req *dto.UserRequest) (*models.UserModel, error) {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.UserModel{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
