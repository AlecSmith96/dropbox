package adapters

import (
	"log/slog"
	"os"
)

type FileWriter struct {
	destinationPath string
}

func NewFileWriter(destinationPath string) *FileWriter {
	return &FileWriter{
		destinationPath: destinationPath,
	}
}

func (writer *FileWriter) CreateFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		slog.Debug("failed to create file", "err", err)
		return err
	}

	// TODO: check if file has any contents to write

	return nil
}

func (writer *FileWriter) DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		slog.Debug("failed to delete file", "err", err)
		return err
	}
	return nil
}
