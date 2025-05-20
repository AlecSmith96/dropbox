run-app:
	go run ./cmd/app/main.go

run-server:
	go run ./cmd/server/main.go

test:
	go test ./pkg/... -coverprofile=coverage.out