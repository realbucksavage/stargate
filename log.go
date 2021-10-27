package stargate

import (
	"fmt"
	"log"
	"os"
)

// Logger is a facade interface that Stargate uses to log its events. By default, Stargate uses a Logger instance that
// writes to os.Stdout with the prefix `STARGATE>`.
type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type defaultLogger struct {
	logger *log.Logger
}

func (d defaultLogger) Info(format string, args ...interface{}) {
	d.logger.Printf("[INF] %s", fmt.Sprintf(format, args...))
}

func (d defaultLogger) Warn(format string, args ...interface{}) {
	d.logger.Printf("[WRN] %s", fmt.Sprintf(format, args...))
}

func (d defaultLogger) Debug(format string, args ...interface{}) {
	d.logger.Printf("[DBG] %s", fmt.Sprintf(format, args...))
}

func (d defaultLogger) Error(format string, args ...interface{}) {
	d.logger.Printf("[ERR] %s", fmt.Sprintf(format, args...))
}

// Log is an instance of Logger used by Stargate.
var Log Logger = defaultLogger{log.New(os.Stdout, "STARGATE> ", log.LstdFlags|log.Lshortfile)}
