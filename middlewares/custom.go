package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		reqID := uuid.New().String()

		// Ajouter un ID unique au header de réponse
		c.Writer.Header().Set("X-Request-ID", reqID)

		// Traiter la requête
		c.Next()

		// Log après exécution
		latency := time.Since(startTime)
		log.Printf("[REQ %s] %s %s | %d | %v",
			reqID,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}
