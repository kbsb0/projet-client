package main

import (
	"ari2-client/controllers"
	"ari2-client/database"
	"ari2-client/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter crée un routeur pour les tests
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/proxy/submit", controllers.SubmitProxyGrid)
	return r
}

func TestSubmitProxyGrid_Validation(t *testing.T) {
	// Initialiser une BDD en mémoire pour le test
	database.Connect() // Note: Idéalement, configurer pour utiliser ":memory:" dans les tests

	r := SetupRouter()

	// Cas 1 : Données invalides (Nom manquant)
	submission := models.Submission{
		Grid: [][]string{}, // Grille vide
	}
	jsonValue, _ := json.Marshal(submission)

	req, _ := http.NewRequest("POST", "/proxy/submit", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// On s'attend à une erreur 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
