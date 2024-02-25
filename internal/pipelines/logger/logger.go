package logger

type Logger interface {
	Printf(format string, v ...any)
	Fatalf(format string, v ...any)
	PrintCommand(char rune)
	Close() error
}
