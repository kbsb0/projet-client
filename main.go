package main

import (
	"ari2-client/controllers"

	"ari2-client/database"

	"ari2-client/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	// 1. Initialiser la BDD

	database.Connect()

	// 2. Configurer Gin (Mode Release ou Debug)

	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 3. Charger les templates HTML

	r.LoadHTMLGlob("templates/*")

	// 4. Appliquer les Middlewares globaux

	r.Use(middlewares.RequestLogger())

	// 5. Définir les routes

	// Route HTML

	r.GET("/", controllers.RenderHome)

	// Routes API (Proxy)

	api := r.Group("/proxy")

	{

		api.GET("/state", controllers.GetProxyState)

		api.POST("/submit", controllers.SubmitProxyGrid)

		// Nouvelle fonctionnalité locale

		api.GET("/history", controllers.GetLocalHistory)

	}

	r.GET("/history", controllers.RenderHistory)

	// 6. Lancer le serveur

	r.Run(":8081")

}
