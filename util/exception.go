package util

import (
	"encoding/json"
	"net/http"
)

func ErrorException(w http.ResponseWriter, err error, errorCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
