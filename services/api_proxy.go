package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Le serveur distant tourne sur le port 8080
const ServerAPI = "http://localhost:8080"

var httpClient = &http.Client{Timeout: 5 * time.Second}

func FetchStateFromRemote() ([]byte, int, error) {
	// Créer la requête GET vers le serveur distant
	req, err := http.NewRequest(http.MethodGet, ServerAPI+"/api/state", nil)
	if err != nil {
		return nil, 0, err
	}

	// Envoyer la requête
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// Retourner corps, status et erreur = nil
	return body, resp.StatusCode, nil
}

func PostGridToRemote(payload any) ([]byte, int, error) {
	jsonData, _ := json.Marshal(payload)
	resp, err := httpClient.Post(ServerAPI+"/api/submit", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}
