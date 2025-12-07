    
package services

import (
	"io"
	"net/http"
	"time"
	"bytes"
	"encoding/json"	
)

const ServerAPI = "http://localhost:8080"

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func FetchStateFromRemote() ([]byte, int, error) {
	resp, err := httpClient.Get(ServerAPI + "/api/state")
	if err != nil {
		return nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
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