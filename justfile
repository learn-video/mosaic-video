clean:
    rm -r hls/*

test:
    go test -v ./...

lint:
    golangci-lint run -v ./...

deps:
    docker compose up

run:
    go run main.go
