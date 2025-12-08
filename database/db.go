package database

import (
	"ari2-client/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("pixel.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Échec de connexion à la BDD:", err)
	}

	// Crée la table submissions automatiquement
	DB.AutoMigrate(&models.Submission{})
	log.Println("Base de données connectée.")
}
