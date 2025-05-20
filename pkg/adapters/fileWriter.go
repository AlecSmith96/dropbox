package adapters

import (
	"log/slog"
	"os"
	"path/filepath"
)

type FileWriter struct {
	destinationPath string
}

func NewFileWriter(destinationPath string) *FileWriter {
	return &FileWriter{
		destinationPath: destinationPath,
	}
}

// CreateFile is a function for creating a file at the given path, with the contents if provided. For any files that exist in sub directories it will
// recursively create each sub directory before creating the file.
func (writer *FileWriter) CreateFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		slog.Debug("failed to create file", "err", err)
		return err
	}

	if len(data) > 0 {
		if _, err := file.Write(data); err != nil {
			slog.Debug("failed to write contents to file", "err", err)
			return err
		}
	}

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
