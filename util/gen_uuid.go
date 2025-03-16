package util

import "github.com/gofrs/uuid"

func GenerateUuid() (uuid.UUID, error) {
	return uuid.NewV4()
}
