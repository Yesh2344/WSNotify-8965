package main

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes the given payload as a JSON response with the provided status code.
// It sets the appropriate Content-Type header and logs any encoding errors.
func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
// small cleanup
		// Use the global logger; this function is deliberately small and shouldn't panic.
		logrus.WithError(err).Error("failed to write JSON response")
	}
}