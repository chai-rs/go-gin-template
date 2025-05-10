package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `gorm:"column:id"`
	Email          string     `gorm:"column:email"`
	HashedPassword string     `gorm:"column:hashed_password"`
	CreatedAt      *time.Time `gorm:"column:created_at"`
}

func (u *User) TableName() string {
	return "users"
}
