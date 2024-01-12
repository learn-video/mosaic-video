test:
    go test -v ./...

lint:
    golangci-lint run -v ./...

deps:
    docker compose up

worker:
    go run main.go worker

storage:
    go run main.go storage
