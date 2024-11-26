package main

import (
	"fmt"
	"log"
	"mori/app"
	"mori/database"
	loginback "mori/loginBack"
	"mori/router"
	"time"
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
	router.InitializeRouter(db)

	// Start a ticker that runs the DeleteNonVerifiedAccounts every 48 hours
	go func() {
		ticker := time.NewTicker(48 * time.Hour) // 48 hours interval
		defer ticker.Stop()

		// Run it initially once at the start
		err := loginback.DeleteNonVerifiedAccounts(db)
		if err != nil {
			log.Printf("Error deleting non-verified accounts: %v", err)
		}

		// Run it every 48 hours
		for range ticker.C {
			err := loginback.DeleteNonVerifiedAccounts(db)
			if err != nil {
				log.Printf("Error deleting non-verified accounts: %v", err)
			}
		}
	}()

	// Pass the router to server.go to start the server
	app.StartServer(router.Router)
}
