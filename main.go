package main

import (
	"ari2-client/controllers"
	"ari2-client/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialiser la BDD
	database.Connect()

	r := gin.Default()

	// Charger les templates
	r.LoadHTMLGlob("templates/*")

	// Route d'accueil
	r.GET("/", controllers.RenderHome)

	// Routes API (Proxy)
	api := r.Group("/proxy")
	{
		api.GET("/state", controllers.GetProxyState)
		api.POST("/submit", controllers.SubmitProxyGrid)
		api.GET("/history", controllers.GetLocalHistory)
	}

	r.Run(":8081")
}