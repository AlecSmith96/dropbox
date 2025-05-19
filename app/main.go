package main

import (
	"context"
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	"log/slog"
	"os"
	"time"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading config", "err", err)
		os.Exit(1)
	}

	_, err = os.Stat(conf.SourceDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("source directory does not exist", "err", err)
			os.Exit(1)
		}
		slog.Error("getting file info for source directory", "err", err)
		os.Exit(1)
	}

	directoryMonitor, err := adapters.NewDirectoryMonitor(conf.SourceDirectory)
	if err != nil {
		slog.Error("creating directory monitor", "err", err)
		os.Exit(1)
	}

	syncEvents := directoryMonitor.SyncDestinationWithSource()
	for _, event := range syncEvents {
		slog.Debug("event", event)
	}

	eventChannel := make(chan entities.FilesystemEvent)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// runs the command line application to poll directory for file changes
	go func() {
		slog.Info("polling for file events")
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err = directoryMonitor.PollForFileChanges(eventChannel)
				if err != nil {
					slog.Error("polling for file changes", "err", err)
					cancel()
				}
			}
		}
	}()

	// process any file events
	for {
		select {
		case event, ok := <-eventChannel:
			if !ok {
				slog.Info("No more events, exiting")
				return
			}

			switch event.Operation {
			case entities.OperationCreated:
				slog.Info("CREATE request", "event", event)
				// CREATE request
			case entities.OperationRenamed:
				slog.Info("RENAMED request", "event", event)
				// UPDATE request
			case entities.OperationDeleted:
				slog.Info("DELETED request", "event", event)
				// DELETE request
			case entities.OperationModified:
				slog.Info("MODIFIED request", "event", event)
				// UPDATE request

			default:
				slog.Error("unknown event operation", "operation", event.Operation)
			}
			slog.Info("processed event", "event", event)
		case <-ctx.Done():
			slog.Debug("shutting down listener")
			return
		}
	}

}
