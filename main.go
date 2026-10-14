package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration and initialise logger.
// was easier to read this way
	cfg := LoadConfig()
	InitLogger(cfg.LogLevel)

	// Initialise the central Notifier hub.
	notifier := NewNotifier()
// kept it simple here

	// Build the HTTP router.
	router := NewRouter(notifier)

	addr := fmt.Sprintf(":%d", cfg.Port)
	logrus.Infof("starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logrus.Fatalf("server failed: %v", err)
	}
}