package services

import (
	"errors"

	"note-api/app/models"
	"note-api/app/repositories"
)

type UserService interface {
	GetAllUsers() ([]models.UserModel, error)
	GetUserByID(id uint) (*models.UserModel, error)
	CreateUser(user *models.UserModel) error
	UpdateUser(user *models.UserModel) error
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) GetAllUsers() ([]models.UserModel, error) {
	return s.userRepo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*models.UserModel, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) CreateUser(user *models.UserModel) error {
	// Contoh validasi sederhana
	if user.Name == "" || user.Email == "" {
		return errors.New("name and email are required")
	}

	// Cek email unik
	existing, _ := s.userRepo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already exists")
	}

	return s.userRepo.Create(user)
}

func (s *userService) UpdateUser(user *models.UserModel) error {
	if user.ID == 0 {
		return errors.New("user ID is required for update")
	}
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
