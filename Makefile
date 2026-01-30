.PHONY: build build-kvit build-kvitd test tidy
build: build-kvit build-kvitd

build-kvit: tidy
	go build -o bin/kvit .

build-kvitd: tidy
	go build -o bin/kvitd ./cmd/kvitd

clean:
	rm -rf bin/

test: tidy
	go test ./...

tidy:
	go mod tidy
