package main

import (
	"ari2-client/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Création du routeur avec les middlewares par défaut (logger + recovery)
	r := gin.Default()

	// Charger les templates HTML
	r.LoadHTMLGlob("templates/*")

	// Route d'accueil utilisant votre contrôleur
	r.GET("/", controllers.RenderHome)

	// Groupe /proxy
	proxy := r.Group("/proxy")
	{
		proxy.GET("/state", controllers.GetProxyState)
	}

	// Lancement du serveur sur le port 8081
	r.Run(":8081")
}
