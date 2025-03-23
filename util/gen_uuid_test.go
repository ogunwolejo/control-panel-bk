package util

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGenerateUuid checks if the function generates a valid UUID.
func TestGenerateUuid(t *testing.T) {
	u, err := GenerateUuid()

	assert.NoError(t, err, "GenerateUuid should not return an error")
	assert.NotEqual(t, uuid.Nil, u, "Generated UUID should not be nil/zero")
}
