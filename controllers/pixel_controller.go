package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ari2-client/services"
	"ari2-client/models"
	"ari2-client/database"
	"encoding/json"
)

func RenderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Pixel Challenge Pro",
	})
}

func GetProxyState(c *gin.Context) {
   data, status, err := services.FetchStateFromRemote()
   if err != nil {
   	c.JSON(status, gin.H{"error": "Serveur API inaccessible"})
   	return
   }
   // On renvoie directement le JSON reçu du serveur distant
   c.Data(status, "application/json", data)
}

func SubmitProxyGrid(c *gin.Context) {
	var submission models.Submission

	//Validation des données entrantes
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Données invalides : " + err.Error(),
		})
		return
	}

	// Logique BDD : Sauvegarde locale de l'historique
	// On convertit la grille en JSON string pour le stockage simple SQLite
	gridBytes, _ := json.Marshal(submission.Grid)
	submission.GridData = string(gridBytes)

	if result := database.DB.Create(&submission); result.Error != nil {
		// On loggue l'erreur mais on ne bloque pas l'envoi au serveur distant
		// C'est un exemple de stratégie "Fail-safe"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde locale"})
		return
	}

	// Appel au serveur distant
	remoteBody, status, err := services.PostGridToRemote(submission)
	if err != nil {
		c.JSON(status, models.APIResponse{Success: false, Message: "Impossible de contacter l'API distante"})
		return
	}

	// Retourne la réponse du serveur distant au client JS
	c.Data(status, "application/json", remoteBody)
}

func GetLocalHistory(c *gin.Context) {
	var history []models.Submission
	// Récupère les 10 derniers envois
	database.DB.Order("created_at desc").Limit(10).Find(&history)

	c.JSON(http.StatusOK, history)
}

// RenderHistory affiche la page d'historique
func RenderHistory(c *gin.Context) {
	c.HTML(http.StatusOK, "history.html", gin.H{
		"title": "Historique des dessins",
	})
}
