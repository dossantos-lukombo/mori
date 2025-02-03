package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mori/pkg/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var accessSecret string
var refreshSecret string

type CustomClaims struct {
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	// History        []struct {
	// 	UserRequest string `json:"user_request"`
	// 	LLMResponse string `json:"llm_response"`
	// } `json:"history"`
	// Stop           bool   `json:"stop"`

	jwt.RegisteredClaims
}

// type Conversation struct {
// 	UserID          string   `json:"user_id"`
// 	ConversationID  string   `json:"conversation_id"`
// 	Session         string   `json:"session"`
// 	UserRequest     string   `json:"user_request"`
// 	LLMResponse     string   `json:"llm_response"`
// 	NewConversation bool     `json:"new_conversation"`
// 	CreatedAt       string   `json:"created_at"`
// 	UpdateAt        string   `json:"update_at"`
// 	History         []string `json:"history"`
// }

type ServerPython struct {
	Status         string `json:"status"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Response       string `json:"response"`
	Timestamp      string `json:"timestamp"`
}

type ConversationResponse struct {
	Type string              `json:"type"`
	Data models.Conversation `json:"data"`
}

func (handler *Handler) LLMHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMHandler"+err.Error(),
				http.StatusInternalServerError)
			return
		}
		fmt.Println("conversation", conversation)

		accessToken, err := GenerateJWT(conversation.UserID, conversation.ConversationID, conversation.UserRequest)
		if err != nil {
			http.Error(w, "Error generating JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating JWT: ", err)
			return

		}

		refreshToken, err := GenerateRefreshJWT(conversation.UserID, conversation.ConversationID, conversation.UserRequest)
		if err != nil {
			http.Error(w, "Error generating refresh JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating refresh JWT: ", err)
			return

		}

		//put token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "accessToken",
			Value:    accessToken,
			Path:     "http://127.0.0.1:8000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 60, // 7 minutes
		})

		// println("JWT Token: ", accessToken)

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Path:     "http://127.0.0.1:8000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 24 * 60 * 60, // 7 jours
		})

		llmConversation := map[string]interface{}{
			"user_id":         conversation.UserID,
			"conversation_id": conversation.ConversationID,
			"message":         conversation.UserRequest,
			// "history":         conversation.History,
			// "stop":            conversation.Stop,
		}

		data, err := json.Marshal(llmConversation)
		if err != nil {
			http.Error(w, "Error marshalling JSON",
				http.StatusInternalServerError)
			return

		}

		SendRequestWithToken("http://127.0.0.1:8000/llm-protected", accessToken, data, w)
		return
	}
}

func SendRequestWithToken(url string, token string, jsonData []byte, w http.ResponseWriter) {

	llm_message := ""
	// Créer une requête POST avec le JSON
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the Authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming non supporté", http.StatusInternalServerError)
		return
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		// Envoyer chaque chunk au frontend
		fmt.Fprintf(w, "%s", line)
		llm_message += string(line)
		// fmt.Println("Response body:", string(line))
		flusher.Flush() // Envoyer immédiatement les données au client
	}

	fmt.Println("Response status stream:", resp.Status)
	// fmt.Println("Response body:", string(responseBody))
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 7)), // Expire dans 1 heure
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

// func SendRequestWithToken(url string, token string, jsonData []byte) {
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		fmt.Println("Erreur lors de la création de la requête :", err)
// 		return
// 	}

// 	// Ajouter l'en-tête Authorization avec le token JWT
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Erreur lors de l'envoi de la requête :", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	fmt.Println("Statut de la réponse :", resp.Status)
// }

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

// LLMConvoSave saves the conversation to the database
func (handler *Handler) LLMConvoSave(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoSave "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		if conversation.NewConversation == true {
			conversation.ConversationID = uuid.NewV4().String()
			fmt.Println("conversation", conversation)
			err = handler.repos.LLMConvoRepo.SaveConvo(conversation)
			if err != nil {
				http.Error(w, "Error saving conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			conversation.ConversationID, err = handler.repos.LLMConvoRepo.GetLastConvoID() //get the last conversation ID
			if err != nil {
				http.Error(w, "Error getting last conversation ID: "+err.Error(), http.StatusInternalServerError)
				return
			}
			err = handler.repos.LLMConvoRepo.SaveConvo(conversation)
			if err != nil {
				http.Error(w, "Error saving conversation: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		return
	}
}

// functon that get all the previous conversation
func (handler *Handler) LLMConvoGetAll(w http.ResponseWriter, convo models.Conversation) {

	var conversations []models.Conversation

	fmt.Println("Convo: ", convo)
	conversations, err := handler.repos.LLMConvoRepo.GetAllConvo(convo)
	if err != nil {
		http.Error(w, "Error getting all conversations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Conversations: ", conversations)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// LLMConvoGet gets the conversation from the database
func (handler *Handler) LLMConvoGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoGet "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		handler.LLMConvoGetAll(w, conversation)

		return
	}
}

func (handler *Handler) LLMConvoGetLast(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request

		type RequestConvo struct {
			UserID string `json:"user_id"`
		}

		var requestConvo RequestConvo

		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body LLMConvoGetLast request: ", string(body))

		err = json.Unmarshal(body, &requestConvo)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoGet "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		conversations, err := handler.repos.LLMConvoRepo.GetLastConvo()
		if err != nil {
			http.Error(w, "Error getting last conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Last conversation: ", conversations)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(conversations)
		if err != nil {
			http.Error(w, "Error encoding JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}

// LLMConvoDelete deletes the conversation from the database
func (handler *Handler) LLMConvoDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation models.Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON LLMConvoDelete "+err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = handler.repos.LLMConvoRepo.DeleteConvo(conversation)
		if err != nil {
			http.Error(w, "Error deleting conversation: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Conversation deleted")

		return
	}
}
