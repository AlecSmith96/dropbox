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
		syncEvents = append(syncEvents, entities.FilesystemEvent{
			Name:         path,
			Operation:    entities.OperationCreated,
			FileContents: metadata,
		})
	}

	return syncEvents
}

func (monitor *DirectoryMonitor) PollForFileChanges(eventChan chan entities.FilesystemEvent) error {
	currentSnapshot, err := monitor.BuildSnapshot(monitor.rootPath)
	if err != nil {
		slog.Error("building snapshot of directory", "err", err)
	}

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

func (monitor *DirectoryMonitor) BuildSnapshot(root string) (map[string]entities.FileContents, error) {
	directoryMap := make(map[string]entities.FileContents)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// log and skip this entry
			slog.Error("error accessing path", "path", path, "err", err)
			return nil
		}

		// skip directories
		if d.IsDir() {
			return nil
		}

		// retrieve FileInfo to get Sys()
		info, err := d.Info()
		if err != nil {
			slog.Error("getting file info", "path", path, "err", err)
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			slog.Error("read file", "path", path, "err", err)
			return nil
		}

		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			slog.Error("unexpected Sys() type", "path", path)
			return nil
		}

		directoryMap[path] = entities.FileContents{
			ModTime: info.ModTime(),
			Size:    info.Size(),
			Inode:   st.Ino,
			Data:    data,
		}
		return nil
	})

	return directoryMap, err
}
