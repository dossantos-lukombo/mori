package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var accessSecret string
var refreshSecret string

type CustomClaims struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`

	jwt.RegisteredClaims
}

// Fonction pour générer un JWT
func GenerateJWT(username, conversationID, message string) (string, error) {
	// Définir les claims
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env : %v", err)
	}
	accessSecret = os.Getenv("ACCESS_SECRET_KEY_LLM")
	claims := CustomClaims{
		UserID:         username,
		ConversationID: conversationID,
		Message:        message,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   username,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingKey := []byte(accessSecret)
	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return tokenString, nil
}

func GenerateRefreshJWT(username, conversationID, message string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Erreur lors du chargement du fichier .env : %v", err)
	}
	refreshSecret = os.Getenv("REFRESH_SECRET_KEY_LLM")

	// Définir les claims
	claims := CustomClaims{
		UserID:         username,
		ConversationID: conversationID,
		Message:        message,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Expire dans 1 heure
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mori",
			Subject:   username,
		},
	}

	// Créer le token avec les claims et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération du token : %v", err)
	}

	return tokenString, nil
}

func SendRequestWithToken(url string, token string, jsonData []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erreur lors de la création de la requête :", err)
		return
	}

	// Ajouter l'en-tête Authorization avec le token JWT
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de l'envoi de la requête :", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Statut de la réponse :", resp.Status)
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Lire le refresh token envoyé par le client
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token manquant", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	// Valider le refresh token
	claims := &CustomClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil {
		http.Error(w, "Refresh token invalide ou expiré", http.StatusUnauthorized)
		return
	}

	// Générer un nouvel access token
	newAccessToken, err := GenerateJWT(claims.UserID, claims.ConversationID, claims.Message)
	if err != nil {
		http.Error(w, "Erreur lors de la génération du nouvel access token", http.StatusInternalServerError)
		return
	}

	// Retourner le nouvel access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": newAccessToken,
	})
}

func VerifyAndRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer les tokens depuis les cookies
	accessTokenCookie, err := r.Cookie("accessToken")
	if err != nil {
		http.Error(w, "Access token manquant", http.StatusUnauthorized)
		return
	}
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token manquant", http.StatusUnauthorized)
		return
	}

	accessToken := accessTokenCookie.Value
	refreshToken := refreshTokenCookie.Value

	// Vérifier l'expiration de l'access token
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})

	// Si le token est valide, vérifier s'il est proche de l'expiration
	if err == nil && token.Valid {
		timeRemaining := time.Until(claims.ExpiresAt.Time)
		if timeRemaining > 2*time.Minute {
			// Token encore valide, pas besoin de le renouveler
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":      "valid",
				"accessToken": accessToken,
			})
			return
		}
	}

	// Si l'access token est expiré ou proche de l'expiration, vérifier le refresh token
	refreshClaims := &CustomClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		http.Error(w, "Refresh token invalide ou expiré", http.StatusUnauthorized)
		return
	}

	// Générer un nouveau access token
	newAccessToken, err := GenerateJWT(refreshClaims.UserID, refreshClaims.ConversationID, refreshClaims.Message)
	if err != nil {
		http.Error(w, "Erreur lors de la génération du nouveau token", http.StatusInternalServerError)
		return
	}

	// Retourner le nouveau access token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "refreshed",
		"accessToken": newAccessToken,
	})
}
