package main

import (
	"fmt"
	"log"
	"mori/app"
	"mori/database"
	config "mori/router"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// Initialize the router and pass the db object
	config.InitializeRouter(db)

	// Pass the router to server.go to start the server
	app.StartServer(config.Router)
}
