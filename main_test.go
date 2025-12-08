package main

import (
	"ari2-client/controllers" // Assurez-vous que le nom du module correspond à votre go.mod
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// SetupRouter permet de créer une instance de routeur isolée pour les tests
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode) // Silence les logs Gin pendant les tests
	r := gin.Default()
	r.POST("/proxy/submit", controllers.SubmitProxyGrid)
	return r
}

func TestSubmitProxyGrid_BadRequest(t *testing.T) {
	// 1. Initialisation
	r := SetupRouter()

	// 2. Préparation de la requête (Body vide "{}" pour provoquer l'erreur de validation)
	req, _ := http.NewRequest("POST", "/proxy/submit", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	// 3. Enregistreur de réponse
	w := httptest.NewRecorder()

	// 4. Exécution
	r.ServeHTTP(w, req)

	// 5. Vérification (CORRECTION ICI)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Code attendu %d (Bad Request), mais reçu %d", http.StatusBadRequest, w.Code)
	}
}