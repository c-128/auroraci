package logger

import (
	"fmt"
)

func NewStdoutLogger() Logger {
	return stdoutLogger(0)
}

type stdoutLogger uint8

func (logger stdoutLogger) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
	fmt.Println()
}

func (logger stdoutLogger) Fatalf(format string, v ...any) {
	fmt.Printf(format, v...)
	fmt.Println()
}

func (logger stdoutLogger) PrintCommand(char rune) {
	fmt.Print(string(char))
}

func (logger stdoutLogger) Close() error {
	return nil
}
