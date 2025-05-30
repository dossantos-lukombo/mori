package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dossantos-lukombo/mori/backend/pkg/models"

	_ "github.com/lib/pq" // PostgreSQL driver
	migrate "github.com/rubenv/sql-migrate"

	"github.com/joho/godotenv"
)

// InitDB initializes the PostgreSQL database connection.
func InitDB() *sql.DB {
	// En dev : charge .env, mais ignore l’erreur si c’est en prod
	_ = godotenv.Load("../.env")

	// 1) Cas Supabase / prod : DATABASE_URL complet
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Failed to connect via DATABASE_URL: %v", err)
		}
		if err := Migrations(db); err != nil {
			log.Fatalf("Failed to apply migrations for prod db: %v", err)
		}
		return db
	}

	// 2) Sinon : mode local Dev
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if host == "" || port == "" || user == "" || password == "" || dbname == "" || sslmode == "" {
		log.Fatal("Missing required environment variables for local Postgres connection.")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to local Postgres: %v", err)
	}
	if err := Migrations(db); err != nil {
		log.Fatalf("Failed to apply migrations for dev db: %v", err)
	}
	return db
}

// InitRepositories initializes all repositories with the database connection.
func InitRepositories(db *sql.DB) *models.Repositories {
	return &models.Repositories{
		UserRepo:     &UserRepository{DB: db},
		SessionRepo:  &SessionRepository{DB: db},
		GroupRepo:    &GroupRepository{DB: db},
		NotifRepo:    &NotifRepository{DB: db},
		MsgRepo:      &MsgRepository{DB: db},
		LLMConvoRepo: &LLMConvoRepository{DB: db},
	}
}

// Migrations applies database migrations.
func Migrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/db/migration/PostgreSql",
	}

	// Execute migrations for PostgreSQL
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Printf("Applied %d migrations to PostgreSQL database!\n", n)
	return nil
}
