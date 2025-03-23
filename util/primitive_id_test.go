package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetPrimitiveID_ValidID tests conversion of a valid hex string to ObjectID
func TestGetPrimitiveID_ValidID(t *testing.T) {
	validID := "5f74d2b4e40e1d4f5c1a2b3c" // A valid 24-character hex string

	objID, err := GetPrimitiveID(validID)

	assert.NoError(t, err, "Valid ObjectID should not return an error")
	assert.NotNil(t, objID, "Valid ObjectID should not be nil")
	assert.Equal(t, validID, objID.Hex(), "Returned ObjectID should match input hex")
}

// TestGetPrimitiveID_InvalidID tests an invalid hex string
func TestGetPrimitiveID_InvalidID(t *testing.T) {
	invalidID := "invalid_hex_string"

	objID, err := GetPrimitiveID(invalidID)

	assert.Error(t, err, "Invalid ObjectID should return an error")
	assert.Nil(t, objID, "Invalid ObjectID should return nil")
}

// TestGetPrimitiveID_InvalidLength tests a string with incorrect length
func TestGetPrimitiveID_InvalidLength(t *testing.T) {
	shortID := "12345" // Too short
	longID := "5f74d2b4e40e1d4f5c1a2b3cabcdef" // Too long

	objID1, err1 := GetPrimitiveID(shortID)
	objID2, err2 := GetPrimitiveID(longID)

	assert.Error(t, err1, "Short ObjectID should return an error")
	assert.Nil(t, objID1, "Short ObjectID should return nil")

	assert.Error(t, err2, "Long ObjectID should return an error")
	assert.Nil(t, objID2, "Long ObjectID should return nil")
}
