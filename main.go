package main

import (
	"ari2-client/controllers"
	"ari2-client/database"
	"ari2-client/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()
	r.Use(middlewares.RequestLogger()) // Votre logger de l'étape 3.1
	r.LoadHTMLGlob("templates/*")

	// --- ROUTES PUBLIQUES (Login/Register) ---
	r.GET("/login", controllers.RenderLogin)
	r.POST("/login", controllers.Login)
	r.GET("/register", controllers.RenderRegister)
	r.POST("/register", controllers.Register)
	r.GET("/logout", controllers.Logout)



	// --- ROUTES PROTÉGÉES (Nécessitent une connexion) ---
	// On crée un groupe qui utilise le middleware d'Auth
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// L'accueil n'est plus public : il faut être loggé pour voir la grille
		protected.GET("/", func(c *gin.Context) {
			// On récupère le nom de l'utilisateur depuis le middleware
			username, _ := c.Get("username")
			c.HTML(200, "index.html", gin.H{
				"title": "Pixel Challenge Pro",
				"username": username, // On l'envoie au template
			})
		})

		protected.GET("/history", controllers.RenderHistory)

		// API Proxy
		api := protected.Group("/proxy")
		{
			api.GET("/state", controllers.GetProxyState)
			api.POST("/submit", controllers.SubmitProxyGrid)
			api.GET("/history", controllers.GetLocalHistory)
			api.POST("/cheat", controllers.CheatHandler)
		}
	}

	r.Run(":8081")
}