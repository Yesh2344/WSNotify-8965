package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// TestNotifierBroadcast verifies that a message broadcast reaches a connected client.
func TestNotifierBroadcast(t *testing.T) {
	// Silence logger output during test.
	logrus.SetOutput(nil)

	// Set up a test HTTP server that upgrades connections to WebSocket.
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade failed: %v", err)
		}
		// Keep the connection open for the test.
		_, _, err = c.NextReader()
		if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			t.Fatalf("client read error: %v", err)
		}
	}))
	defer srv.Close()

	// Convert http://127.0.0.1/... to ws://127.0.0.1/...
	wsURL := "ws" + srv.URL[4:] // strip "http"

	// Dial the WebSocket endpoint.
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	defer ws.Close()

	// Initialise Notifier and register the test client.
	notifier := NewNotifier()
	notifier.AddClient(ws)

	// Broadcast a test message.
	testMsg := "hello test"
	notifier.Broadcast([]byte(testMsg))

	// Verify the client receives the message.
	_, msg, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if string(msg) != testMsg {
		t.Fatalf("expected %q, got %q", testMsg, string(msg))
	}
}