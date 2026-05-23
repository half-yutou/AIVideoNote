.PHONY: run build test clean

run:
	go run ./cmd/server/ -config config.yaml -log-level debug

build:
	go build -o bin/aivideonote.exe ./cmd/server/

test:
	go test ./...

clean:
	rm -rf bin/ data/
