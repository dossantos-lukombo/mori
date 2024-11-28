package loginback

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// DeleteNonVerifiedAccounts deletes users who have not verified their email within 48 hours
func DeleteNonVerifiedAccounts(db *sql.DB) error {
	// Calculate the time 48 hours ago
	thresholdTime := time.Now().Add(-48 * time.Hour) // Reverting back to 48 hours

	// Query to find users who have not verified their email and whose account creation date is older than 48 hours
	query := `
		DELETE FROM users
		WHERE verified = FALSE
		AND create_date < $1
		RETURNING id, username, email;`

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(thresholdTime)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Log the users being deleted
	var deletedCount int
	for rows.Next() {
		var id uint
		var username, email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("Deleted user: ID = %d, Username = %s, Email = %s", id, username, email)
		deletedCount++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %v", err)
	}

	// If no rows were deleted, log that as well
	if deletedCount == 0 {
		log.Println("No non-verified users were deleted.")
	}

	return nil
}
