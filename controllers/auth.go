package controllers

import (
	"ari2-client/database"
	"ari2-client/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Clé secrète pour signer les tokens (à mettre dans une var d'env idéalement)
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func init() {
    // Valeur par défaut pour le TP si pas de variable d'env
    if len(jwtKey) == 0 {
        jwtKey = []byte("super_secret_key_tp_gin")
    }
}

// 1. Inscription
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

    // Hashage du mot de passe
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de cryptage"})
		return
	}

	user := models.User{Username: username, Password: string(hash)}

    // Sauvegarde en BDD
    if result := database.DB.Create(&user); result.Error != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{"error": "Ce pseudo est déjà pris !"})
		return
	}

	// Succès -> vers login
	c.Redirect(http.StatusFound, "/login")
}

// 2. Connexion
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Utilisateur inconnu"})
		return
	}

    // Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Mot de passe incorrect"})
		return
	}

	// Génération du Token JWT
	expirationTime := time.Now().Add(1 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur génération token"})
		return
	}

    // Stockage du token dans un Cookie sécurisé (HttpOnly)
	c.SetCookie("auth_token", tokenString, 3600, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/")
}

// 3. Déconnexion
func Logout(c *gin.Context) {
    // On écrase le cookie avec une date passée
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// Affichage des pages
func RenderLogin(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) }
func RenderRegister(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) }