package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Hash a password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword"), bcrypt.DefaultCost)

	// Successful Login Test Case
	t.Run("Successful Login", func(t *testing.T) {
		// Mock database query
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "session", "verified"}).
			AddRow(1, "testuser", "testuser@example.com", hashedPassword, "testsession", true)
		mock.ExpectQuery("SELECT id, username, email, password, session, verified FROM users WHERE").
			WithArgs("testuser", "testuser").
			WillReturnRows(rows)

		// Create a request
		req := httptest.NewRequest("POST", "/login", strings.NewReader("username_email=testuser&password=securepassword"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// Call the handler
		handler := LoginHandler(db)
		handler.ServeHTTP(rr, req)

		// Assert response
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %v", rr.Code)
		}
		if rr.Body.String() != `{"message": "Login successful"}`+"\n" {
			t.Errorf("Unexpected response body: %v", rr.Body.String())
		}
	})

	// Incorrect Credentials Test Case
	t.Run("Incorrect Credentials", func(t *testing.T) {
		// Mock database query returning no rows
		mock.ExpectQuery("SELECT id, username, email, password, session, verified FROM users WHERE").
			WithArgs("wronguser", "wronguser").
			WillReturnError(sql.ErrNoRows)

		// Create a request
		req := httptest.NewRequest("POST", "/login", strings.NewReader("username_email=wronguser&password=wrongpassword"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// Call the handler
		handler := LoginHandler(db)
		handler.ServeHTTP(rr, req)

		// Assert response
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code 401, got %v", rr.Code)
		}
		if rr.Body.String() != `{"error": "Invalid username/email or password"}`+"\n" {
			t.Errorf("Unexpected response body: %v", rr.Body.String())
		}
	})

	// Unverified User Test Case
	t.Run("Unverified User", func(t *testing.T) {
		// Mock database query
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "session", "verified"}).
			AddRow(1, "unverifieduser", "unverifieduser@example.com", hashedPassword, "testsession", false)
		mock.ExpectQuery("SELECT id, username, email, password, session, verified FROM users WHERE").
			WithArgs("unverifieduser", "unverifieduser").
			WillReturnRows(rows)

		// Create a request
		req := httptest.NewRequest("POST", "/login", strings.NewReader("username_email=unverifieduser&password=securepassword"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// Call the handler
		handler := LoginHandler(db)
		handler.ServeHTTP(rr, req)

		// Assert response
		if rr.Code != http.StatusForbidden {
			t.Errorf("Expected status code 403, got %v", rr.Code)
		}
		if rr.Body.String() != `{"error": "Please verify your email before logging in"}`+"\n" {
			t.Errorf("Unexpected response body: %v", rr.Body.String())
		}
	})

	// Internal Server Error Test Case
	t.Run("Internal Server Error", func(t *testing.T) {
		// Mock database query returning an error
		mock.ExpectQuery("SELECT id, username, email, password, session, verified FROM users WHERE").
			WithArgs("erroruser", "erroruser").
			WillReturnError(sql.ErrConnDone)

		// Create a request
		req := httptest.NewRequest("POST", "/login", strings.NewReader("username_email=erroruser&password=securepassword"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// Call the handler
		handler := LoginHandler(db)
		handler.ServeHTTP(rr, req)

		// Assert response
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %v", rr.Code)
		}
		if rr.Body.String() != `{"error": "Internal server error"}`+"\n" {
			t.Errorf("Unexpected response body: %v", rr.Body.String())
		}
	})
}
