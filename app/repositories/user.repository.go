package repositories

import (
	"gorm.io/gorm"

	"{{.ModuleName}}/app/models"
)

type UserRepository interface {
	FindAll() ([]models.UserModel, error)
	FindByID(id uint) (*models.UserModel, error)
	FindByEmail(email string) (*models.UserModel, error)
	Create(user *models.UserModel) error
	Update(user *models.UserModel) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll() ([]models.UserModel, error) {
	var users []models.UserModel
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByID(id uint) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.UserModel) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.UserModel) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserModel{}, id).Error
}
