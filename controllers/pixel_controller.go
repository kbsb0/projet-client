package controllers

import (
	"ari2-client/database"
	"ari2-client/models"
	"ari2-client/services"
	"encoding/json"
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

func SubmitProxyGrid(c *gin.Context) {
	var submission models.Submission

	// 1. Validation (Binding)
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Données invalides : " + err.Error(),
		})
		return
	}

	// 2. Préparation BDD (Conversion Grid -> String)
	gridBytes, _ := json.Marshal(submission.Grid)
	submission.GridData = string(gridBytes)

	// 3. Sauvegarde locale
	if result := database.DB.Create(&submission); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde locale"})
		return
	}

	// 4. Envoi au serveur distant
	remoteBody, status, err := services.PostGridToRemote(submission)
	if err != nil {
		// En cas d'erreur réseau, on prévient le client mais la donnée est sauvée en local
		c.JSON(status, models.APIResponse{Success: false, Message: "Impossible de contacter l'API distante"})
		return
	}

	c.Data(status, "application/json", remoteBody)
}

func GetLocalHistory(c *gin.Context) {
	var history []models.Submission
	// Récupère les 10 derniers, du plus récent au plus ancien
	database.DB.Order("created_at desc").Limit(10).Find(&history)
	c.JSON(http.StatusOK, history)
}

func RenderHistory(c *gin.Context) {
	c.HTML(http.StatusOK, "history.html", gin.H{
		"title": "Historique Local",
	})
}
