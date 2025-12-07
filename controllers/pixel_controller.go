package controllers

import (
	"ari2-client/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RenderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Pixel Challenge Pro",
	})
}

func GetProxyState(c *gin.Context) {
	body, status, err := services.FetchStateFromRemote()

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Impossible de contacter le serveur distant",
		})
		return
	}

	c.Data(status, "application/json", body)
}
