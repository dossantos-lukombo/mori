package main

import (
	"fmt"
	"io"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../frontend/chatbot/chatbot.html")

		if r.Method == "POST" {
			// Do something

			//get the body of our POST request

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			bodyString := string(body)
			fmt.Println(bodyString)

		}

	}
}
