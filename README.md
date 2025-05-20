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

### Config
The valid config values for this application are:

| Name                    | Example Values                              | Description                                                                           |
|-------------------------|---------------------------------------------|---------------------------------------------------------------------------------------|
| `source-directory`      | ./some/relative/path OR /some/absolute/path | The location of the source directory to monitor.                                      |
| `destination-directory` | ./some/relative/path OR /some/absolute/path | The location of the destination directory to synchronise with the source directory.   |
| `base-url`              | http://localhost:8080                       | The base URL the server is running on.                                                |
| `use-absolute-paths`    | true OR false                               | Whether to use relative or absolute paths for the source and destination directories. |
NOTE: When using absolute paths, `~` is not expanded by Go, so if you want to use a path like `~/some/absolute/path`, then the application will need to be passed `/some/absolute/path` with the config value `use-absolute-paths` set to `true`.

For the purpose of this exercise, a `dev-config.yaml` file has been committed in the repository with the config variables
already added to allow for quick setup. The source and destination directories must exist for the application to start, if not you will get
an error similar to:
```
2025/05/19 21:51:24 ERROR source directory does not exist err="stat ./tmp/different/src: no such file or directory"
```

As the`app` synchronises initial source state on startup, the `server` needs to be ready when it starts up. However, when the app starts up it will poll a liveness endpoint on the server until it is ready to serve traffic before beginning the sync. This means that the app and server can be started in any order.

### Running the tests
All the code written in this application was written with testing in mind. It is logically grouped into structs that represent a function in the code, using dependency injection
to allow for easy mocking of dependencies.

To run the unit tests you can use this command:
```
go test ./pkg/... -coverprofile=coverage.out
```

The tests utilise `go.uber.org/mock` to generate mocks from interfaces. These can then force errors in the tests to increase coverage. 
Due to time constraints, I haven't been able to hit the test coverage I would normally like to, but I would use unit tests in conjunction with integration tests as well
to ensure that any technologies used with the application work correctly when communicating with it (such as databases or message topics).

## Assumptions

During the development of the application I have made a few assumptions, these are:

- On startup, the destination directory will always be empty or reflect the source directory exactly, and not require synchronising the directories 
bi-directionally. Any differences in the destination directory will not be updated. The `app` will synchronise any existing files in the source directory on startup.