Étape 1.2 : Serveur Web & HTML

Dans pixel_controller :
```
package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RenderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Pixel Challenge Pro",
	})
}

```

Dans main.go :
```

package main

import (
	"ari2-client/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Charger les templates
	r.LoadHTMLGlob("templates/*")

	// Route d'accueil
	r.GET("/", controllers.RenderHome)

	r.Run(":8081")
}

```






Étape 1.3 : La couche Service & Proxy simple

Fichier : services/api_proxy.go
code Go
```
    
package services

import (
	"io"
	"net/http"
	"time"
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

  ```

Mise à jour : controllers/pixel_controller.go
code Go

 ```   
import "ari2-client/services" // Ajouter l'import

// ... (RenderHome existant)

func GetProxyState(c *gin.Context) {
	data, status, err := services.FetchStateFromRemote()
	if err != nil {
		c.JSON(status, gin.H{"error": "Serveur API inaccessible"})
		return
	}
	// On renvoie directement le JSON reçu du serveur distant
	c.Data(status, "application/json", data)
}

  ```

Mise à jour : main.go
code Go
```
    
// ...
	// Routes API (Proxy)
	api := r.Group("/proxy")
	{
		api.GET("/state", controllers.GetProxyState)
	}
// ...

  ```




Étape 1.3 : La couche Service & Proxy simple

Fichier : services/api_proxy.go
code Go

```
    
package services

import (
	"io"
	"net/http"
	"time"
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

  ```

Mise à jour : controllers/pixel_controller.go
code Go
```
    
import "ari2-client/services" // Ajouter l'import

// ... (RenderHome existant)

func GetProxyState(c *gin.Context) {
	data, status, err := services.FetchStateFromRemote()
	if err != nil {
		c.JSON(status, gin.H{"error": "Serveur API inaccessible"})
		return
	}
	// On renvoie directement le JSON reçu du serveur distant
	c.Data(status, "application/json", data)
}
```

```
Mise à jour : main.go
code Go

    
// ...
	// Routes API (Proxy)
	api := r.Group("/proxy")
	{
		api.GET("/state", controllers.GetProxyState)
	}
// ...

  ```