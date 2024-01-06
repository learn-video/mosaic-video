clean:
    rm -r output/*

test:
    go test -v ./...

deps:
    docker compose up

run:
    go run main.go
