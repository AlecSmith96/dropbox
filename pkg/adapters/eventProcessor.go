package adapters

import (
	"errors"
	"fmt"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	"log/slog"
	"strings"
)

type EventProcessor struct {
	httpClient *HTTPClient
	sourcePath string
}

func NewEventProcessor(client *HTTPClient, sourcePath string) *EventProcessor {
	var path string
	trimmedSourcePath := strings.Split(sourcePath, "./")
	if len(trimmedSourcePath) != 2 {
		path = sourcePath
	} else {
		path = trimmedSourcePath[1]
	}

	return &EventProcessor{
		httpClient: client,
		sourcePath: path,
	}
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
		err := processor.httpClient.SendCreateRequest(filePathWithoutSource, event.FileContents.Data, event.FileContents.IsDirectory)
		if err != nil {
			slog.Error("processing create request", "err", err)
		}
	case entities.OperationRenamed:
		// UPDATE request
	case entities.OperationDeleted:
		err := processor.httpClient.SendDeleteRequest(filePathWithoutSource)
		if err != nil {
			slog.Error("processing create request", "err", err)
		}
	case entities.OperationModified:
		// UPDATE request

	default:
		slog.Error("unknown event operation", "operation", event.Operation)
		return fmt.Errorf("unknown event operation: %s", event.Operation)
	}

	return nil
}
