package main

import (
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/drivers"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading config", "err", err)
		os.Exit(1)
	}

	destinationDirectory := conf.DestinationDirectory
	if conf.UseAbsolutePaths {
		home, _ := os.UserHomeDir()
		destinationDirectory = filepath.Join(home, conf.DestinationDirectory[1:])
	}

	// validate the destination directory exists
	_, err = os.Stat(destinationDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("destination directory does not exist", "err", err, "destinationDir", destinationDirectory)
			os.Exit(1)
		}
		slog.Error("getting file info for destination directory", "err", err)
		os.Exit(1)
	}

	// init dependencies
	fileWriter := adapters.NewFileWriter(destinationDirectory)
	router := drivers.NewRouter(destinationDirectory, fileWriter)

	router.Run()
}
