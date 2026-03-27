.PHONY: dev build run clean

dev:
	air

build:
	go build -o ./tmp/main ./cmd

run:
	go run ./cmd

clean:
	rm -rf ./tmp
