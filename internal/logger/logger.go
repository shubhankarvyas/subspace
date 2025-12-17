package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log severity
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel Level = INFO
	logger       *log.Logger
)

// logEntry represents a structured log entry
type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Init initializes the logger with the specified level
func Init(level string) {
	logger = log.New(os.Stdout, "", 0)
	
	switch level {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		currentLevel = INFO
	}
}

// Debug logs a debug message with optional key-value pairs
func Debug(msg string, keysAndValues ...interface{}) {
	if currentLevel <= DEBUG {
		writeLog("DEBUG", msg, keysAndValues...)
	}
}

// Info logs an info message with optional key-value pairs
func Info(msg string, keysAndValues ...interface{}) {
	if currentLevel <= INFO {
		writeLog("INFO", msg, keysAndValues...)
	}
}

// Warn logs a warning message with optional key-value pairs
func Warn(msg string, keysAndValues ...interface{}) {
	if currentLevel <= WARN {
		writeLog("WARN", msg, keysAndValues...)
	}
}

// Error logs an error message with optional key-value pairs
func Error(msg string, keysAndValues ...interface{}) {
	if currentLevel <= ERROR {
		writeLog("ERROR", msg, keysAndValues...)
	}
}

// writeLog is the internal logging function that handles structured output
func writeLog(level, msg string, keysAndValues ...interface{}) {
	if logger == nil {
		Init("info")
	}

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]interface{}),
	}

	// Parse key-value pairs
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := fmt.Sprint(keysAndValues[i])
			entry.Fields[key] = keysAndValues[i+1]
		}
	}

	// Output as JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("Failed to marshal log entry: %v", err)
		return
	}

	logger.Println(string(jsonData))
}

// WithContext creates a contextual logger that automatically adds fields to all logs
type ContextLogger struct {
	module string
	fields map[string]interface{}
}

// NewContext creates a new contextual logger
func NewContext(module string, fields ...interface{}) *ContextLogger {
	cl := &ContextLogger{
		module: module,
		fields: make(map[string]interface{}),
	}
	
	cl.fields["module"] = module
	
	// Parse additional fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprint(fields[i])
			cl.fields[key] = fields[i+1]
		}
	}
	
	return cl
}

// Debug logs with context
func (cl *ContextLogger) Debug(msg string, keysAndValues ...interface{}) {
	Debug(msg, cl.mergeFields(keysAndValues...)...)
}

// Info logs with context
func (cl *ContextLogger) Info(msg string, keysAndValues ...interface{}) {
	Info(msg, cl.mergeFields(keysAndValues...)...)
}

// Warn logs with context
func (cl *ContextLogger) Warn(msg string, keysAndValues ...interface{}) {
	Warn(msg, cl.mergeFields(keysAndValues...)...)
}

// Error logs with context
func (cl *ContextLogger) Error(msg string, keysAndValues ...interface{}) {
	Error(msg, cl.mergeFields(keysAndValues...)...)
}

// mergeFields combines context fields with new fields
func (cl *ContextLogger) mergeFields(keysAndValues ...interface{}) []interface{} {
	result := make([]interface{}, 0, len(cl.fields)*2+len(keysAndValues))
	
	// Add context fields first
	for k, v := range cl.fields {
		result = append(result, k, v)
	}
	
	// Add new fields
	result = append(result, keysAndValues...)
	
	return result
}

// Timing logs the duration of an operation
func Timing(module, action string, start time.Time, err error) {
	duration := time.Since(start)
	fields := []interface{}{
		"module", module,
		"action", action,
		"duration_ms", duration.Milliseconds(),
	}
	
	if err != nil {
		fields = append(fields, "error", err.Error())
		Error("Action completed with error", fields...)
	} else {
		Info("Action completed", fields...)
	}
}
