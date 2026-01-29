.PHONY: build test tidy
build: tidy
	go build -o kvit .

test: tidy
	go test ./...

tidy:
	go mod tidy
