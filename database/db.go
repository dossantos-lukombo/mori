package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// User struct for application-level logic
type User struct {
	ID          uint
	Username    string
	Email       string
	Password    string
	Session     string
	FavorisJSON string
}

// Conversation struct for application-level logic
type Conversation struct {
	ID       uint
	UserID   uint
	Title    string
	Echanges string
	Dates    string
}

// InitDB initializes the database connection and creates tables
func InitDB() (*sql.DB, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Printf("Database is unreachable: %v", err)
		return nil, err
	}

	// Create tables
	err = createTables(db)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
		return nil, err
	}

	log.Println("Database initialized and tables created successfully!")
	return db, nil
}

// createTables runs raw SQL statements to create the required tables
func createTables(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		session TEXT NOT NULL,
		favoris_json JSONB DEFAULT '[]'
	);`

	conversationTable := `
	CREATE TABLE IF NOT EXISTS conversations (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		title TEXT NOT NULL,
		echanges TEXT,
		dates TEXT
	);`

	// Execute the SQL statements
	_, err := db.Exec(userTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	_, err = db.Exec(conversationTable)
	if err != nil {
		return fmt.Errorf("failed to create conversations table: %v", err)
	}

	return nil
}

// UpdateFavoris updates the FavorisJSON field for a user
func UpdateFavoris(db *sql.DB, userID uint, newFavoris []uint) error {
	// Convert favoris slice to JSON
	favorisJSON, err := json.Marshal(newFavoris)
	if err != nil {
		return err
	}

	query := `UPDATE users SET favoris_json = $1 WHERE id = $2;`
	_, err = db.Exec(query, string(favorisJSON), userID)
	return err
}

// GetFavoris retrieves the FavorisJSON field and converts it to a []uint
func GetFavoris(db *sql.DB, userID uint) ([]uint, error) {
	query := `SELECT favoris_json FROM users WHERE id = $1;`
	var favorisJSON string
	err := db.QueryRow(query, userID).Scan(&favorisJSON)
	if err != nil {
		return nil, err
	}

	// Parse JSON into []uint
	var favoris []uint
	err = json.Unmarshal([]byte(favorisJSON), &favoris)
	if err != nil {
		return nil, err
	}

	return favoris, nil
}
