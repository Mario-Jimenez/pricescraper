package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitializeLogger set the initial configuration fot the global app logger
func InitializeLogger(serviceName, serviceVersion, logLevel string) {
	// Log as JSON
	log.SetFormatter(utcFormatter{&log.JSONFormatter{}, serviceName, serviceVersion})

	// Output to file
	log.SetOutput(os.Stdout)

	// Only log specified severity or above
	switch logLevel {
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warning":
		log.SetLevel(log.WarnLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}

	// Add the calling method as a field
	log.SetReportCaller(true)
}
