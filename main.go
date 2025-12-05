package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ServerAPI  = "http://localhost:8080" // Adresse du vrai serveur
	ClientPort = ":8081"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// On sert juste la page, en injectant l'URL de l'API Serveur
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"ServerAPI": ServerAPI})
	})

	r.Run(ClientPort)
}
