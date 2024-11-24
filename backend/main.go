package main

import (
	"fmt"
	"log"
	"mori/app"
	"mori/database"
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
	InitializeRouter(db)

	// Pass the router to server.go to start the server
	app.StartServer(Router)
}
