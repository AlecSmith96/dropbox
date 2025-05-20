package main

import (
	"github.com/AlecSmith96/dopbox/pkg/adapters"
	"github.com/AlecSmith96/dopbox/pkg/drivers"
	"log/slog"
	"os"
)

func main() {
	conf, err := adapters.NewConfig()
	if err != nil {
		slog.Error("reading config", "err", err)
		os.Exit(1)
	}

	_, err = os.Stat(conf.DestinationDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("destination directory does not exist", "err", err, "destinationDir", conf.DestinationDirectory)
			os.Exit(1)
		}
		slog.Error("getting file info for destination directory", "err", err)
		os.Exit(1)
	}

	fileWriter := adapters.NewFileWriter(conf.DestinationDirectory)
	router := drivers.NewRouter(conf.DestinationDirectory, fileWriter)

	router.Run()
}
