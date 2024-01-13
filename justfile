test:
    go test -v -race ./...

lint:
    golangci-lint run -v ./...

deps:
    docker compose up

worker:
    go run main.go worker

storage:
    go run main.go storage

player:
    go run main.go player
