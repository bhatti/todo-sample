.PHONY: build test lint clean

build:
	go build ./...

test:
	go test -race ./...

lint:
	go vet ./...

clean:
	rm -rf bin/
