package main

import (
	"context"
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/entities"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading config", "err", err)
		os.Exit(1)
	}

	sourceDirectory := conf.SourceDirectory
	if conf.UseAbsolutePaths {
		home, _ := os.UserHomeDir()
		sourceDirectory = filepath.Join(home, conf.SourceDirectory[1:])
	}

	// validate the source directory exists
	_, err = os.Stat(sourceDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("source directory does not exist", "err", err)
			os.Exit(1)
		}
		slog.Error("getting file info for source directory", "err", err)
		os.Exit(1)
	}

	// init dependencies
	httpClient := adapters.NewHTTPClient(http.DefaultClient, conf.BaseURL)
	serverLive := false
	slog.Info("checking sever liveness")
	for !serverLive {
		serverLive = httpClient.IsServerLive()
	}
	slog.Info("server live!")

	directoryMonitor, err := adapters.NewDirectoryMonitor(sourceDirectory)
	if err != nil {
		slog.Error("creating directory monitor", "err", err)
		os.Exit(1)
	}

	eventProcessor := adapters.NewEventProcessor(httpClient, sourceDirectory)

	// make sure all files in source directory on startup get created in destination
	syncEvents := directoryMonitor.SyncDestinationWithSource()
	for _, event := range syncEvents {
		err := eventProcessor.ProcessEvent(event)
		if err != nil {
			slog.Error("processing event", "err", err)
			continue
		}
	}

	if len(syncEvents) > 1 {
		slog.Info("synced existing files in destination")
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

			err := eventProcessor.ProcessEvent(event)
			if err != nil {
				slog.Error("processing event", "err", err)
				continue
			}

			slog.Info("processed filesystem event", "operation", event.Operation, "filePath", event.Name)
		case <-ctx.Done():
			slog.Debug("shutting down listener")
			return
		}
	}

}
