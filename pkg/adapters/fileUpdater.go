package adapters

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type FileUpdater struct {
}

func NewFileUpdater() *FileUpdater {
	return &FileUpdater{}
}

// SyncDestinationWithSource is a function to synchronise the entire source directory with the destination directory.
func (updater *FileUpdater) SyncDestinationWithSource(source, destination string) error {
	sourceDirInfo, err := os.Stat(source)
	if err != nil {
		slog.Debug("unable to open source directory", "err", err)
		return err
	}

	destDirInfo, err := os.Stat(destination)
	if err != nil {
		slog.Debug("unable to open destination directory", "err", err)
		return err
	}

	if !sourceDirInfo.IsDir() || !destDirInfo.IsDir() {
		return errors.New("source and destination paths must be directories")
	}

	updater.CopyDirectoryContents(destination)

	return nil
}

func (updater *FileUpdater) CopyDirectoryContents(filePath string) error {
	err := filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If WalkDir failed to read a directory, report and skip it.
			slog.Debug("skipping", "path", path, "err", err)
			return fs.SkipDir
		}

		if d.IsDir() {
			// TODO: recursively call static function
		}

		base := d.Name()
		nameOnly := strings.TrimSuffix(base, filepath.Ext(base))
		destDir := filepath.Join(filePath, nameOnly)

		// create the directory (and parents) if necessary
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %q: %w", destDir, err)
		}

		// open source file
		srcFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open src %q: %w", path, err)
		}
		defer srcFile.Close()

		// create destination file
		dstPath := filepath.Join(destDir, base)
		dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("create dst %q: %w", dstPath, err)
		}
		defer dstFile.Close()

		// copy contents
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return fmt.Errorf("copy %q to %q: %w", path, dstPath, err)
		}

		return nil
	})
	if err != nil {
		slog.Error("copying directory contents", "err", err)
		return err
	}

	return nil
}
