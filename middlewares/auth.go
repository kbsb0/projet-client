package middlewares

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super_secret_key_tp_gin") // Même clé que dans controllers

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
        // 1. Récupérer le cookie
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			// Pas de cookie ? -> Login
            // Si c'est une requête AJAX (API), on renvoie 401
            // Si c'est une requête navigateur (Page), on redirige
            if c.Request.Header.Get("Accept") == "application/json" {
			    c.AbortWithStatus(http.StatusUnauthorized)
            } else {
                c.Redirect(http.StatusFound, "/login")
                c.Abort()
            }
			return
		}

        // 2. Parser et valider le token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
            c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

        // 3. Extraire les données (Claims) et les mettre dans le contexte Gin
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // C'est ici qu'on stocke le "username" pour l'utiliser dans les controlleurs
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}