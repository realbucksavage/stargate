package stargate

import (
	log "github.com/op/go-logging"
)

var (
	loggerName = "stargate"

	// Logger logs all events in Stargate
	Logger = log.MustGetLogger(loggerName)
)

// SetLevel changes the level of Stargate logger.
func SetLevel(level log.Level) {
	log.SetLevel(level, loggerName)
}
