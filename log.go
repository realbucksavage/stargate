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

// Info logs strings formatted by `format` on INFO level.
func (d defaultLogger) Info(format string, args ...interface{}) {
	d.logger.Printf("[INF] %s", fmt.Sprintf(format, args...))
}

// Warn logs strings formatted by `format` on WARNING level.
func (d defaultLogger) Warn(format string, args ...interface{}) {
	d.logger.Printf("[WRN] %s", fmt.Sprintf(format, args...))
}

// Debug logs strings formatted by `format` on DEBUG level.
func (d defaultLogger) Debug(format string, args ...interface{}) {
	d.logger.Printf("[DBG] %s", fmt.Sprintf(format, args...))
}

// Error logs strings formatted by `format` on ERROR level.
func (d defaultLogger) Error(format string, args ...interface{}) {
	d.logger.Printf("[ERR] %s", fmt.Sprintf(format, args...))
}

// Log is an instance of Logger used by Stargate. It is set to an implementation of stargate.Logger that writes to
// the standard output. Implementors may update this variable to their own implementation of stargate.Logger.
var Log Logger = defaultLogger{log.New(os.Stdout, "STARGATE> ", log.LstdFlags|log.Lshortfile)}
