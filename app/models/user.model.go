package models

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        string    `json:"id" gorm:"type:char(26);primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
}

func (UserModel) TableName() string {
	return "users"
}

func (user *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		t := time.Now()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		user.ID = ulid.MustNew(ulid.Timestamp(t), entropy).String()
	}
	return
}
