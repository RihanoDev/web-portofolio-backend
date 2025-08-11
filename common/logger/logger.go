package logger

import (
"fmt"
"io"
"os"
"path/filepath"

"github.com/sirupsen/logrus"
)

// LogLevel represents log levels
type LogLevel int

const (
DebugLevel LogLevel = iota
InfoLevel
WarnLevel
ErrorLevel
FatalLevel
)

// Fields represents structured logging fields
type Fields map[string]interface{}

// Logger interface defines logging contract
type Logger interface {
Debug(message string, fields ...Fields)
Info(message string, fields ...Fields)
Warn(message string, fields ...Fields)
Error(message string, fields ...Fields)
Fatal(message string, fields ...Fields)
WithFields(fields Fields) Logger
}

// AppLogger implements Logger interface using logrus
type AppLogger struct {
logger *logrus.Logger
fields logrus.Fields
}

// NewLogger creates a new application logger
func NewLogger(level LogLevel, outputs ...io.Writer) Logger {
log := logrus.New()

// Set log level
switch level {
case DebugLevel:
log.SetLevel(logrus.DebugLevel)
case InfoLevel:
log.SetLevel(logrus.InfoLevel)
case WarnLevel:
log.SetLevel(logrus.WarnLevel)
case ErrorLevel:
log.SetLevel(logrus.ErrorLevel)
case FatalLevel:
log.SetLevel(logrus.FatalLevel)
default:
log.SetLevel(logrus.InfoLevel)
}

// Set formatter
log.SetFormatter(&logrus.JSONFormatter{
TimestampFormat: "2006-01-02 15:04:05",
})

// Set output
if len(outputs) > 0 {
log.SetOutput(io.MultiWriter(outputs...))
} else {
// Default: write to both file and stdout
logFile := createLogFile()
if logFile != nil {
log.SetOutput(io.MultiWriter(os.Stdout, logFile))
} else {
log.SetOutput(os.Stdout)
}
}

return &AppLogger{
logger: log,
fields: make(logrus.Fields),
}
}

// NewDevelopmentLogger creates a logger for development environment
func NewDevelopmentLogger() Logger {
return NewLogger(DebugLevel)
}

// NewProductionLogger creates a logger for production environment
func NewProductionLogger() Logger {
return NewLogger(InfoLevel)
}

// Debug logs a debug message
func (l *AppLogger) Debug(message string, fields ...Fields) {
l.logWithFields(logrus.DebugLevel, message, fields...)
}

// Info logs an info message
func (l *AppLogger) Info(message string, fields ...Fields) {
l.logWithFields(logrus.InfoLevel, message, fields...)
}

// Warn logs a warning message
func (l *AppLogger) Warn(message string, fields ...Fields) {
l.logWithFields(logrus.WarnLevel, message, fields...)
}

// Error logs an error message
func (l *AppLogger) Error(message string, fields ...Fields) {
l.logWithFields(logrus.ErrorLevel, message, fields...)
}

// Fatal logs a fatal message and exits
func (l *AppLogger) Fatal(message string, fields ...Fields) {
l.logWithFields(logrus.FatalLevel, message, fields...)
}

// WithFields creates a new logger with additional fields
func (l *AppLogger) WithFields(fields Fields) Logger {
newFields := make(logrus.Fields)
// Copy existing fields
for k, v := range l.fields {
newFields[k] = v
}
// Add new fields
for k, v := range fields {
newFields[k] = v
}

return &AppLogger{
logger: l.logger,
fields: newFields,
}
}

// logWithFields logs a message with optional fields
func (l *AppLogger) logWithFields(level logrus.Level, message string, fields ...Fields) {
entry := l.logger.WithFields(l.fields)

// Add any additional fields
for _, fieldSet := range fields {
for k, v := range fieldSet {
entry = entry.WithField(k, v)
}
}

entry.Log(level, message)
}

// createLogFile creates a log file for writing
func createLogFile() *os.File {
logDir := "log"
if err := os.MkdirAll(logDir, 0755); err != nil {
fmt.Printf("Failed to create log directory: %v\n", err)
return nil
}

logFile, err := os.OpenFile(
filepath.Join(logDir, "app.log"),
os.O_CREATE|os.O_WRONLY|os.O_APPEND,
0666,
)
if err != nil {
fmt.Printf("Failed to open log file: %v\n", err)
return nil
}

return logFile
}

// RequestLogger creates a logger with request context
type RequestLogger struct {
Logger
requestID string
}

// NewRequestLogger creates a logger with request context
func NewRequestLogger(baseLogger Logger, requestID string) *RequestLogger {
return &RequestLogger{
Logger:    baseLogger.WithFields(Fields{"request_id": requestID}),
requestID: requestID,
}
}

// LogRequest logs HTTP request details
func (rl *RequestLogger) LogRequest(method, path, clientIP string, duration int64) {
rl.Info("HTTP request", Fields{
"method":    method,
"path":      path,
"client_ip": clientIP,
"duration":  duration,
})
}

// LogError logs an error with request context
func (rl *RequestLogger) LogError(err error, context string) {
rl.Error(fmt.Sprintf("Error in %s: %v", context, err), Fields{
"error":   err.Error(),
"context": context,
})
}

// LogBusinessOperation logs business operation
func (rl *RequestLogger) LogBusinessOperation(operation string, entity string, entityID interface{}) {
rl.Info("Business operation", Fields{
"operation": operation,
"entity":    entity,
"entity_id": entityID,
})
}
