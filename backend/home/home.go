package home

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mori/middleware"
	"net/http"
	"path/filepath"
)

type Conversation struct {
	UserID          string   `json:"user_id"`
	ConversationID  string   `json:"conversation_id"`
	Session         string   `json:"session"`
	UserRequest     string   `json:"user_request"`
	LLMResponse     string   `json:"llm_response"`
	NewConversation bool     `json:"new_conversation"`
	History         []string `json:"history"`
}

type ServerPython struct {
	Status         string `json:"status"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Response       string `json:"response"`
	Timestamp      string `json:"timestamp"`
}

type ConversationResponse struct {
	Type string       `json:"type"`
	Data Conversation `json:"data"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method: ", r.Method)
	if r.Method == http.MethodGet {
		// w.Header().Set("Content-Type", "text/html")
		// w.WriteHeader(http.StatusOK)
		path, err := filepath.Abs("../frontend/chatbot/chatbot.html")
		if err != nil {
			http.Error(w, "Error resolving file path", http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, r, path)
		return

	}

	if r.Method == http.MethodPost {
		// Do something
		//get the body of our POST request
		var conversation Conversation
		w.Header().Set("Content-Type", "application/json")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		fmt.Println("Body: ", string(body))

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON",
				http.StatusInternalServerError)
			return
		}
		fmt.Println("conversation", conversation)

		accessToken, err := middleware.GenerateJWT(conversation.UserID, conversation.ConversationID, conversation.UserRequest)
		if err != nil {
			http.Error(w, "Error generating JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating JWT: ", err)
			return

		}

		refreshToken, err := middleware.GenerateRefreshJWT(conversation.UserID, conversation.ConversationID, conversation.UserRequest)
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
			Path:     "http://localhost:8000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 60, // 7 minutes
		})

		// println("JWT Token: ", accessToken)

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    refreshToken,
			Path:     "http://localhost:8000/llm-protected",
			HttpOnly: true,
			Secure:   true, // Activez HTTPS en production
			SameSite: http.SameSiteStrictMode,
			MaxAge:   7 * 24 * 60 * 60, // 7 jours
		})

		llmConversation := map[string]interface{}{
			"user_id":         conversation.UserID,
			"conversation_id": conversation.ConversationID,
			"message":         conversation.UserRequest,
		}

		data, err := json.Marshal(llmConversation)
		if err != nil {
			http.Error(w, "Error marshalling JSON",
				http.StatusInternalServerError)
			return

		}

		// data := []byte(`{"user_id":"` + conversation.UserID + `","conversation_id":"` + conversation.ConversationID + `","message":"` + conversation.UserRequest + `"}`)

		SendRequestWithToken("http://localhost:8000/llm-protected", accessToken, data, w)
		return
	}
}

func SendRequestWithToken(url string, token string, jsonData []byte, w http.ResponseWriter) {

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
	// w.Header().Set("Connection", "keep-alive")

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
		// fmt.Println("Response body:", string(line))
		flusher.Flush() // Envoyer immédiatement les données au client
	}

	fmt.Println("Response status stream:", resp.Status)
	// fmt.Println("Response body:", string(responseBody))
}
