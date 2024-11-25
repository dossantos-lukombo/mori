package main

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

		// fmt.Fprintf(w, "GET request received")
		// Do something

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

		err = json.Unmarshal(body, &conversation)
		if err != nil {
			http.Error(w, "Error unmarshalling JSON",
				http.StatusInternalServerError)
			return
		}
		fmt.Println("conversation", conversation)

		jwtToken, err := middleware.GenerateJWT(conversation.UserID, conversation.ConversationID, conversation.UserRequest)
		if err != nil {
			http.Error(w, "Error generating JWT",
				http.StatusInternalServerError)
			fmt.Println("Error generating JWT: ", err)
			return

		}

		println("JWT Token: ", jwtToken)

		data := []byte(`{"user_id":"` + conversation.UserID + `","conversation_id":"` + conversation.ConversationID + `","message":"` + conversation.UserRequest + `"}`)

		SendRequestWithToken("http://localhost:8000/llm-protected", jwtToken, data, w)

	}
}

func SendRequestWithToken(url string, token string, jsonData []byte, w http.ResponseWriter) {
	// var data []byte
	// var serverPython ServerPython
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the Authorization header with the JWT token
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	http.Error(w, "Error reading response body", http.StatusInternalServerError)
	// 	return
	// }
	// defer resp.Body.Close()
	// fmt.Println("Response body:", string(body))
	// flusher, ok := w.(http.Flusher)
	// if !ok {
	// 	http.Error(w, "Streaming non supporté", http.StatusInternalServerError)
	// 	return
	// }
	// _, err = resp.Body.Read(data)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return
	// }
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
		flusher.Flush() // Envoyer immédiatement les données au client
	}

	fmt.Println("Response status:", resp.Status)
	// fmt.Println("Response body:", string(responseBody))
}
