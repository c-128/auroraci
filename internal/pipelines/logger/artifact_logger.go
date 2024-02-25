package logger

import (
	"bytes"
	"fmt"

	"github.com/c-128/auroraci/internal/pipelines/artifacts"
)

func NewArtifactLogger(filename string, uploader artifacts.Uploader) Logger {
	logger := &artifactLogger{
		filename: filename,
		uploader: uploader,
	}
	return logger
}

type artifactLogger struct {
	filename string
	uploader artifacts.Uploader
	buffer   string
}

func (logger *artifactLogger) Printf(format string, v ...any) {
	logger.buffer += fmt.Sprintf(format, v...) + "\n"
}

func (logger *artifactLogger) Fatalf(format string, v ...any) {
	logger.buffer += fmt.Sprintf(format, v...) + "\n"
}

func (logger *artifactLogger) PrintCommand(char rune) {
	logger.buffer += string(char)
}

func (logger *artifactLogger) Close() error {
	err := logger.uploader.Upload(
		logger.filename,
		bytes.NewBufferString(logger.buffer),
	)
	if err != nil {
		return err
	}

	return nil
}
