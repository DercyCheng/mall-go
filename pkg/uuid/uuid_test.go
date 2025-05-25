package uuid

import (
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	// Generate a new UUID
	id := New()

	// Verify that it's a valid UUID
	if !IsValid(id) {
		t.Errorf("New() generated an invalid UUID: %s", id)
	}

	// Verify that calling New() again generates a different UUID
	id2 := New()
	if id == id2 {
		t.Errorf("New() generated the same UUID twice: %s", id)
	}
}

func TestIsValid(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    string
		expected bool
	}{
		{uuid.New().String(), true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"invalid-uuid", false},
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"", false},
		{"123e4567-e89b-12d3-a456-42661417400", false},   // Too short
		{"123e4567-e89b-12d3-a456-4266141740000", false}, // Too long
	}

	// Run test cases
	for _, tc := range testCases {
		result := IsValid(tc.input)
		if result != tc.expected {
			t.Errorf("IsValid(%s) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestParse(t *testing.T) {
	// Valid UUID
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	parsed, err := Parse(validUUID)
	if err != nil {
		t.Errorf("Parse(%s) returned an error: %v", validUUID, err)
	}
	if parsed.String() != validUUID {
		t.Errorf("Parse(%s) = %s, expected %s", validUUID, parsed.String(), validUUID)
	}

	// Invalid UUID
	invalidUUID := "invalid-uuid"
	_, err = Parse(invalidUUID)
	if err == nil {
		t.Errorf("Parse(%s) did not return an error as expected", invalidUUID)
	}
}

func TestMustParse(t *testing.T) {
	// Valid UUID
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	parsed := MustParse(validUUID)
	if parsed.String() != validUUID {
		t.Errorf("MustParse(%s) = %s, expected %s", validUUID, parsed.String(), validUUID)
	}

	// Invalid UUID should panic, so we recover
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParse with invalid UUID did not panic as expected")
		}
	}()

	// This should panic
	MustParse("invalid-uuid")
}

func TestNil(t *testing.T) {
	nilUUID := Nil()
	expectedNil := "00000000-0000-0000-0000-000000000000"

	if nilUUID != expectedNil {
		t.Errorf("Nil() = %s, expected %s", nilUUID, expectedNil)
	}

	// Verify it's a valid UUID
	if !IsValid(nilUUID) {
		t.Errorf("Nil() returned an invalid UUID: %s", nilUUID)
	}
}
