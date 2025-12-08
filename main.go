package main

import (
	"ari2-client/controllers"
	"ari2-client/database"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	// Création du routeur avec les middlewares par défaut (logger + recovery)
	r := gin.Default()

	// Charger les templates HTML
	r.LoadHTMLGlob("templates/*")

	// Route d'accueil utilisant votre contrôleur
	r.GET("/", controllers.RenderHome)

	// Routes API
	api := r.Group("/proxy")
	{
		api.GET("/state", controllers.GetProxyState)
		api.POST("/submit", controllers.SubmitProxyGrid) // <--- AJOUTER ICI
		api.GET("/history", controllers.GetLocalHistory) // <--- AJOUTER ICI
	}

	r.GET("/history", controllers.RenderHistory) // <--- AJOUTER ICI

	// Lancement du serveur sur le port 8081
	r.Run(":8081")
}
