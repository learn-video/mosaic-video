clean:
    rm -r output/*

test:
    go test -v ./...

lint:
    golangci-lint run -v ./...

deps:
    docker compose up

run:
    go run main.go
