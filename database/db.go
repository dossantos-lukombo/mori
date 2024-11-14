package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey"`
	Username      string         `gorm:"unique;not null"`
	Email         string         `gorm:"unique;not null"`
	Password      string         `gorm:"not null"`
	Session       string         `gorm:"not null"`
	Conversations []Conversation `gorm:"foreignKey:UserID"`
	Favoris       []uint         `gorm:"type:int[]"`
}

type Conversation struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`
	Title    string `gorm:"not null"`
	Echanges string `gorm:"type:text"`
	Dates    string `gorm:"type:text"`
}

func InitDB() (*gorm.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Build the DSN using environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&User{}, &Conversation{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
