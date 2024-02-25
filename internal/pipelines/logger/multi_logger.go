package logger

import "errors"

func NewMultiLogger(logger ...Logger) Logger {
	return multiLogger(logger)
}

type multiLogger []Logger

func (logger multiLogger) Printf(format string, v ...any) {
	for _, logger := range logger {
		logger.Printf(format, v...)
	}
}

func (logger multiLogger) Fatalf(format string, v ...any) {
	for _, logger := range logger {
		logger.Fatalf(format, v...)
	}
}

func (logger multiLogger) PrintCommand(char rune) {
	for _, logger := range logger {
		logger.PrintCommand(char)
	}
}

func (logger multiLogger) Close() error {
	closeErrors := make([]error, len(logger))
	for i, logger := range logger {
		closeErrors[i] = logger.Close()
	}

	err := errors.Join(closeErrors...)
	if err != nil {
		return err
	}

	return nil
}
