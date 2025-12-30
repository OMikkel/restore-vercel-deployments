package logger

import "fmt"

const (
	LevelDebug = iota
	LevelInfo
	LevelError
	LevelDisabled
)

type Logger struct {
	Enabled bool
	Level   int
}

func NewLogger(level int) *Logger {
	return &Logger{
		Level: level,
	}
}

func (l *Logger) Info(messages ...any) {
	if l.Level <= LevelInfo {
		fmt.Println(append([]any{"[INFO]"}, messages...)...)
	}
}

func (l *Logger) Error(messages ...any) {
	if l.Level <= LevelError {
		fmt.Println(append([]any{"[ERROR]"}, messages...)...)
	}
}

func (l *Logger) Debug(messages ...any) {
	if l.Level <= LevelDebug {
		fmt.Println(append([]any{"[DEBUG]"}, messages...)...)
	}
}
