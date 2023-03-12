.PHONY: build clean test

build:
	go build -o build/vischeck ./cmd/vischeck
	go build -buildmode=plugin ./plugin

clean:
	rm -rf build

test:
	go test -coverpkg=./... -coverprofile=/tmp/vischeck.coverage.out ./...
	go tool cover -func=/tmp/vischeck.coverage.out
