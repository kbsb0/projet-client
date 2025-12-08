package controllers

import (
	"ari2-client/database"
	"ari2-client/models"
	"ari2-client/services"
	"encoding/json"
	"net/http"
	"sync"

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



func CheatHandler(c *gin.Context) {
	// 1. Récupérer l'état distant pour avoir la solution
	body, _, err := services.FetchStateFromRemote()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Impossible de lire l'état distant"})
		return
	}

	// 2. Décoder le JSON reçu
	var state models.RemoteState
	if err := json.Unmarshal(body, &state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de décodage JSON"})
		return
	}

	// 3. Transformer la TargetGrid (0/1) en Grid de couleurs (string)
	// On triche en mettant tout en bleu (#3498db) là où il y a un 1.
	var cheatGrid [][]string
	for _, row := range state.TargetGrid {
		var colorRow []string
		for _, cell := range row {
			if cell == 1 {
				colorRow = append(colorRow, "#3498db") // Bleu
			} else {
				colorRow = append(colorRow, "") // Blanc/Vide
			}
		}
		cheatGrid = append(cheatGrid, colorRow)
	}

	// Récupération du pseudo connecté
	username, _ := c.Get("username")

	// Préparation de l'objet à envoyer
	submission := models.Submission{
		Name: username.(string), // Assertion de type
		Grid: cheatGrid,
	}

	// --- LOGIQUE CONCURRENTE ---

	var wg sync.WaitGroup // Création du compteur
	nbWorkers := 5        // Nombre d'envois simultanés

	// On lance 5 goroutines
	for i := 0; i < nbWorkers; i++ {
		wg.Add(1) // On incrémente le compteur AVANT de lancer la goroutine

		go func(workerID int) {
			defer wg.Done() // On décrémente quand la fonction se termine

			// Envoi de la requête (on ignore les erreurs ici pour simplifier le TP)
			services.PostGridToRemote(submission)

			// Optionnel : un petit log pour voir que c'est parallèle
			// fmt.Printf("Worker %d a fini son envoi\n", workerID)
		}(i)
	}

	// Bloque l'exécution ici tant que le compteur du WaitGroup n'est pas à 0
	wg.Wait()

	// 4. Réponse au client
	c.JSON(http.StatusOK, gin.H{
		"message": "💥 C'est fait ! 5 grilles parfaites envoyées.",
		"success": true,
	})
}
