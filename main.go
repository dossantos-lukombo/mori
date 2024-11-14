package main

import (
	"fmt"
	"log"
	app "mori/app"
	db "mori/database" // Update this import path to match your project structure
)

func main() {
	app.StartServer()
	// Initialize the database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the database connection
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve generic database object: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	// Example usage: Create a new user
	user := db.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "securepassword",
		Session:  "initial_session_token",
	}

	result := database.Create(&user)
	if result.Error != nil {
		log.Fatalf("Failed to create user: %v", result.Error)
	}

	fmt.Printf("User %s created successfully!\n", user.Username)
}
