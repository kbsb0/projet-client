package database

import (
	"ari2-client/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	// Crée un fichier local "pixel.db"
	DB, err = gorm.Open(sqlite.Open("pixel.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Échec de connexion à la base de données:", err)
	}

	// Migration automatique du schéma
	DB.AutoMigrate(&models.Submission{})
	log.Println("Base de données connectée et migrée.")
}
