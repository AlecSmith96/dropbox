package adapters

import (
	"fmt"
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
func (writer *FileWriter) CreateFile(path string, data []byte, isDirectory bool) error {
	if isDirectory {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("mkdir %q: %w", path, err)
		}
		return nil
	}

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
	// RemoveAll also handles any sub content, so if the file is a directory with files within it, all files within it
	// are removed as well
	err := os.RemoveAll(path)
	if err != nil {
		slog.Debug("failed to delete file", "err", err)
		return err
	}
	return nil
}

func (writer *FileWriter) RenameFile(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		slog.Debug("unable to rename file", "err", err)
		return err
	}

	return nil
}

func (writer *FileWriter) UpdateFile(path string, data []byte) error {
	// WriteFile will remove all content in the file, then replace it with the new content
	err := os.WriteFile(path, data, 0o644)
	if err != nil {
		slog.Debug("failed to update files contents", "err", err)
		return err
	}

	return nil
}
