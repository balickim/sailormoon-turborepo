package utils

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

func HandleDBError(err error) error {
	if err == nil {
		return nil
	}

	// Handle SQLite-specific errors
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
		// SQLite does not always provide detailed field info, so we may need to generalize
		return fmt.Errorf("field is already taken")
	}

	// Handle PostgreSQL-specific errors
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		// Extract the column name from the PostgreSQL error detail (if available)
		field := extractPostgresField(pqErr)
		if field == "" {
			field = "field" // Fallback in case we cannot determine the field
		}
		return fmt.Errorf("field %s is already taken", field)
	}

	// General error fallback
	return err
}

// extractPostgresField extracts the field name from the PostgreSQL error, if available
func extractPostgresField(err *pq.Error) string {
	// The Detail field of pq.Error sometimes contains info like "Key (email)=(example@example.com) already exists."
	if err.Detail != "" {
		return parseFieldFromDetail(err.Detail)
	}
	return ""
}

// parseFieldFromDetail parses the field name from the PostgreSQL error detail string
func parseFieldFromDetail(detail string) string {
	// Example detail string: "Key (email)=(example@example.com) already exists."
	// We want to extract the field name "email"
	start := "Key ("
	end := ")="

	startIdx := strings.Index(detail, start)
	endIdx := strings.Index(detail, end)

	if startIdx != -1 && endIdx != -1 && startIdx < endIdx {
		return detail[startIdx+len(start) : endIdx]
	}

	return ""
}
