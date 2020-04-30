package stargate

import (
	log "github.com/op/go-logging"
)

var (
	loggerName = "stargate"

	Logger = log.MustGetLogger(loggerName)
)

func SetLevel(level log.Level) {
	log.SetLevel(level, loggerName)
}
