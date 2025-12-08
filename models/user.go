package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"` // Le pseudo doit être unique
	Password string // Contiendra le hash bcrypt
}