package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// wsUpgrader configures the parameters for upgrading HTTP connections to WebSocket.
var wsUpgrader = websocket.Upgrader{
	// Allow any origin for demo purposes; tighten in production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleWebSocket upgrades the HTTP request, registers the client with the Notifier,
// and keeps the connection alive until an error or close occurs.
func handleWebSocket(notifier *Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.WithError(err).Error("WebSocket upgrade failed")
			http.Error(w, "Could not upgrade to websocket", http.StatusBadRequest)
			return
		}
		notifier.AddClient(conn)
		defer func() {
			notifier.RemoveClient(conn)
			_ = conn.Close()
		}()

		// Simple read loop to detect client disconnects.
		for {
			if _, _, err := conn.NextReader(); err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					logrus.Info("WebSocket closed normally")
				} else {
					logrus.WithError(err).Warn("WebSocket read error")
				}
				break
			}
		}
	}
}

// handleNotify receives a JSON payload with a "message" field and broadcasts it.
func handleNotify(notifier *Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
			return
		}
		if payload.Message == "" {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "message cannot be empty"})
			return
		}
		notifier.Broadcast([]byte(payload.Message))
		WriteJSON(w, http.StatusOK, map[string]string{"status": "sent"})
	}
}

// NewRouter assembles the HTTP routes for the server.
func NewRouter(notifier *Notifier) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWebSocket(notifier))
	mux.HandleFunc("/notify", handleNotify(notifier))
	return mux
}