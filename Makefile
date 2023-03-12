.PHONY: build clean lint test

build:
	go build -o build/vischeck ./cmd/vischeck
	go build -buildmode=plugin ./plugin

clean:
	rm -rf build

lint:
	go vet ./...

test:
	go test -coverpkg=./... -coverprofile=/tmp/vischeck.coverage.out ./...
	go tool cover -func=/tmp/vischeck.coverage.out
