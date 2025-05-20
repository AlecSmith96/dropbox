package adapters

import (
	"bytes"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"syscall"
)

type DirectoryMonitor struct {
	rootPath         string
	previousSnapshot map[string]entities.FileContents
}

func NewDirectoryMonitor(root string) (*DirectoryMonitor, error) {
	monitor := &DirectoryMonitor{
		rootPath: root,
	}

	initialSnapshot, err := monitor.BuildSnapshot(root)
	if err != nil {
		return nil, err
	}
	monitor.previousSnapshot = initialSnapshot

	return monitor, nil
}

// SyncDestinationWithSource is a function that returns all files in the filesystem as CREATED events to be sent to
// insert them in the destination directory. This assumes the destination directory starts out as an empty dir, although
// any existing files will be ignored. (Existing files could cause unexpected errors such as a file already exists).
func (monitor *DirectoryMonitor) SyncDestinationWithSource() []entities.FilesystemEvent {
	syncEvents := make([]entities.FilesystemEvent, 0)
	for path, metadata := range monitor.previousSnapshot {
		if path == monitor.rootPath {
			continue
		}

		syncEvents = append(syncEvents, entities.FilesystemEvent{
			Name:         path,
			Operation:    entities.OperationCreated,
			FileContents: metadata,
		})
	}

	return syncEvents
}

// PollForFileChanges is a function that build a snapshot of the current state of the source directory and compares it
// with the previous snapshot. Any differences between them are then calculated and an entities.FilesystemEvent is put on
// the event channel.
func (monitor *DirectoryMonitor) PollForFileChanges(eventChan chan entities.FilesystemEvent) error {
	currentSnapshot, err := monitor.BuildSnapshot(monitor.rootPath)
	if err != nil {
		slog.Error("building snapshot of directory", "err", err)
	}

	// use inodes to compare files as they dont change, allows for tracking filename changes
	previousFilepathByInodes := make(map[uint64]string, len(monitor.previousSnapshot))
	for path, metadata := range monitor.previousSnapshot {
		previousFilepathByInodes[metadata.Inode] = path
	}

	currentFilepathByInodes := make(map[uint64]string, len(currentSnapshot))
	for path, metadata := range currentSnapshot {
		currentFilepathByInodes[metadata.Inode] = path
	}

	for path, metadata := range currentSnapshot {
		oldPath, exists := previousFilepathByInodes[metadata.Inode]
		if !exists {
			eventChan <- entities.FilesystemEvent{
				Name:         path,
				Operation:    entities.OperationCreated,
				FileContents: metadata,
			}
			continue
		}

		if oldPath != path {
			eventChan <- entities.FilesystemEvent{
				Name:         path,
				Operation:    entities.OperationRenamed,
				PreviousPath: oldPath,
				FileContents: metadata,
			}
			continue
		}

		previousFileVersion, exists := monitor.previousSnapshot[path]
		if exists {
			if !bytes.Equal(metadata.Data, previousFileVersion.Data) {
				eventChan <- entities.FilesystemEvent{
					Name:         path,
					Operation:    entities.OperationModified,
					FileContents: metadata,
				}
			}
		}
	}

	for path, metadata := range monitor.previousSnapshot {
		_, exists := currentFilepathByInodes[metadata.Inode]
		if !exists {
			eventChan <- entities.FilesystemEvent{
				Name:         path,
				Operation:    entities.OperationDeleted,
				FileContents: metadata,
			}
		}
	}

	monitor.previousSnapshot = currentSnapshot
	return nil
}

// BuildSnapshot is a function that uses the built in filepath.WalkDir function to traverse the root directory and build
// a map of filepaths to entities.FileContents. It stores whether each file is a directory, its contents and its inode.
func (monitor *DirectoryMonitor) BuildSnapshot(root string) (map[string]entities.FileContents, error) {
	directoryMap := make(map[string]entities.FileContents)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("error accessing path", "path", path, "err", err)
			return nil
		}

		// exclude source directory entry
		if d.Name() == path {
			return nil
		}

		// retrieve FileInfo to get iNode of file
		info, err := d.Info()
		if err != nil {
			slog.Error("getting file info", "path", path, "err", err)
			return err
		}

		// get the iNode of the file to track name changes later on
		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			slog.Error("unexpected Sys() type", "path", path)
			return err
		}

		// dont try to get content if file is a directory
		if d.IsDir() {
			directoryMap[path] = entities.FileContents{
				Inode:       st.Ino,
				IsDirectory: true,
			}
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file", "path", path, "err", err)
			return err
		}

		directoryMap[path] = entities.FileContents{
			Inode:       st.Ino,
			Data:        data,
			IsDirectory: false,
		}
		return nil
	})

	return directoryMap, err
}
