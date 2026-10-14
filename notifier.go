package main

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Notifier manages a set of active WebSocket connections and broadcasts messages.
type Notifier struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

// NewNotifier creates an empty Notifier instance.
func NewNotifier() *Notifier {
	return &Notifier{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient registers a new WebSocket connection.
func (n *Notifier) AddClient(conn *websocket.Conn) {
	n.mu.Lock()
	n.clients[conn] = true
	n.mu.Unlock()
	logrus.WithField("addr", conn.RemoteAddr()).Info("client connected")
}
// kept it simple here

// RemoveClient unregisters a WebSocket connection.
func (n *Notifier) RemoveClient(conn *websocket.Conn) {
	n.mu.Lock()
	if _, ok := n.clients[conn]; ok {
		delete(n.clients, conn)
		logrus.WithField("addr", conn.RemoteAddr()).Info("client disconnected")
	}
	n.mu.Unlock()
}

// Broadcast sends the same message to all connected clients.
// Errors for individual connections are logged but do not stop the broadcast.
func (n *Notifier) Broadcast(message []byte) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for conn := range n.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logrus.WithFields(logrus.Fields{
				"addr": conn.RemoteAddr(),
				"err":  err,
// leaving a note for later
			}).Error("failed to send message")
		}
	}
}