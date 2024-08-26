package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper for zap.Logger
type Logger struct {
	*zap.Logger
}

var (
	instance *Logger
	once     sync.Once
)

// GetInstance returns the singleton instance of Logger
func GetInstance() *Logger {
	once.Do(func() {
		// Create a file to store logs
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		// Create a core that writes logs to the file
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey:    "message",
				LevelKey:      "level",
				TimeKey:       "time",
				CallerKey:     "caller",
				StacktraceKey: "stacktrace",
			}),
			zapcore.AddSync(file),
			zapcore.DebugLevel,
		)

		// Create a logger with the core
		logger := zap.New(core)
		instance = &Logger{logger}
	})

	return instance
}

// Sync flushes the logger
func (l *Logger) Sync() {
	l.Logger.Sync()
}
