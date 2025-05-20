# Dropbox

A simple directory synchronisation tool that monitors a source directory and replicates any updates to a destination directory.

It is made up of two parts:

- `server`: An HTTP server that can perform CRUD operations on the destination directory.
- `app`: A command line application that monitors the source directory and sends update requests to the server via HTTP.

## Getting Started
To run the application you need to run each component separately from the project root:
```
go run ./cmd/app/main.go
```

```
go run ./cmd/server/main.go
```
This is to ensure that the config file can be found by the config adapter.

For the purpose of this exercise, a `dev-config.yaml` file has been committed in the repository with the config variables
already added to allow for quick setup. The source and destination directories must exist for the application to start, if not you will get
an error similar to:
```
2025/05/19 21:51:24 ERROR source directory does not exist err="stat ./tmp/different/src: no such file or directory"
```

## Assumptions

During the development of the application I have made a few assumptions, these are:

- On startup, the destination directory will always be empty, and not require synchronising the directories 
bi-directionally. The `app` will synchronise any existing files in the source directory on startup.