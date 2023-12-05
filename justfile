clean:
    rm -r output/*

test:
    go test -v ./...

run:
    docker compose up
