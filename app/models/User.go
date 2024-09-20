package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint       `gorm:"primaryKey;column:id" json:"id"`
	Name      string     `gorm:"size:256;column:name" json:"name" binding:"required"`
	Email     string     `gorm:"size:256;column:email;not null;unique" json:"email" binding:"omitempty,email"`
	Password  string     `gorm:"size:100;not null" json:"password" binding:"required"`
	Status    bool       `gorm:"default:true" json:"status"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
