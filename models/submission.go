package models

import (
	"gorm.io/gorm"
	"time"
)

// Submission représente un envoi de grille
// Elle sert à la fois pour le binding JSON et la BDD SQLite
type Submission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name" binding:"required"`
	GridData  string         `json:"-" gorm:"type:text"` // On stocke la grille en string dans la BDD
	Grid      [][]string     `json:"grid" gorm:"-"`      // Ignoré par GORM, utilisé par JSON
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// APIResponse structure standard pour tes réponses JSON
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
