// Package uuid provides utilities for generating and handling UUIDs
package uuid

import (
	"github.com/google/uuid"
)

// New generates a new random UUID
func New() string {
	return uuid.New().String()
}

// NewV4 generates a new random UUID (alias for New)
func NewV4() string {
	return New()
}

// IsValid checks if a string is a valid UUID
func IsValid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// Parse parses a UUID string
func Parse(u string) (uuid.UUID, error) {
	return uuid.Parse(u)
}

// MustParse parses a UUID string and panics on error
func MustParse(u string) uuid.UUID {
	return uuid.MustParse(u)
}

// Nil returns the nil UUID (all zeros)
func Nil() string {
	return uuid.Nil.String()
}
