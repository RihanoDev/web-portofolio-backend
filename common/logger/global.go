package logger

// defaultLogger holds the app-wide logger used when no specific logger is injected.
var defaultLogger Logger = NewDevelopmentLogger()

// SetDefaultLogger sets the global default logger instance.
func SetDefaultLogger(l Logger) {
	if l != nil {
		defaultLogger = l
	}
}

// GetLogger returns the global default logger instance.
func GetLogger() Logger {
	return defaultLogger
}
