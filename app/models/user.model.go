package models

import "time"

type UserModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
}

func (UserModel) TableName() string {
	return "users"
}
