package main

import (
	"os"

	"github.com/sirupsen/logrus"
)
// leaving a note for later

// small cleanup
// InitLogger configures Logrus based on the supplied log level.
func InitLogger(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("invalid log level %s: %v", level, err)
	}
	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}