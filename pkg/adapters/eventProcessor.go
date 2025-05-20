package adapters

import (
	"errors"
	"fmt"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	"log/slog"
	"strings"
)

type EventProcessor struct {
	requestSender RequestSender
	sourcePath    string
}

func NewEventProcessor(requestSender RequestSender, sourcePath string) *EventProcessor {
	var path string
	trimmedSourcePath := strings.Split(sourcePath, "./")
	if len(trimmedSourcePath) != 2 {
		path = sourcePath
	} else {
		path = trimmedSourcePath[1]
	}

	return &EventProcessor{
		requestSender: requestSender,
		sourcePath:    path,
	}
}

// RequestSender is an interface that is used to mock the request client calls for testing.
//
//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/requestSender.go  . "RequestSender"
type RequestSender interface {
	SendCreateRequest(path string, data []byte, isDirectory bool) error
	SendDeleteRequest(path string) error
	SendRenameRequest(oldPath, newPath string) error
	SendUpdateRequest(path string, data []byte) error
}

// ProcessEvent is a function that takes a filesystem event and sends the appropriate request to the http server to
// replicate it in the destination directory.
// If a message fails to send, it will log an error message. ideally it would be put on a queue to be processed later
// and prevent data loss.
func (processor *EventProcessor) ProcessEvent(event entities.FilesystemEvent) error {
	trimmedPath := strings.Split(event.Name, processor.sourcePath)
	if len(trimmedPath) != 2 {
		return errors.New("invalid trimmed path produced")
	}
	filePathWithoutSource := trimmedPath[1]

	switch event.Operation {
	case entities.OperationCreated:
		err := processor.requestSender.SendCreateRequest(filePathWithoutSource, event.FileContents.Data, event.FileContents.IsDirectory)
		if err != nil {
			slog.Error("processing create request", "err", err)
			return err
		}

	case entities.OperationRenamed:
		trimmedPath = strings.Split(event.PreviousPath, processor.sourcePath)
		if len(trimmedPath) != 2 {
			return errors.New("invalid trimmed path produced")
		}
		previousFilePathWithoutSource := trimmedPath[1]
		err := processor.requestSender.SendRenameRequest(previousFilePathWithoutSource, filePathWithoutSource)
		if err != nil {
			slog.Error("processing create request", "err", err)
			return err
		}

	case entities.OperationDeleted:
		err := processor.requestSender.SendDeleteRequest(filePathWithoutSource)
		if err != nil {
			slog.Error("processing create request", "err", err)
			return err
		}

	case entities.OperationModified:
		err := processor.requestSender.SendUpdateRequest(filePathWithoutSource, event.FileContents.Data)
		if err != nil {
			slog.Error("processing create request", "err", err)
			return err
		}

	default:
		slog.Error("unknown event operation", "operation", event.Operation)
		return fmt.Errorf("unknown event operation: %s", event.Operation)
	}

	return nil
}
