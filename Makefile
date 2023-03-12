.PHONY: build clean lint test

build:
	go build -o build/vischeck ./cmd/vischeck
	go build -buildmode=plugin -o build/vischeck.so ./plugin

clean:
	rm -rf build tmp

lint:
	go vet ./...

test:
	mkdir -p tmp
	go test -coverpkg=./... -coverprofile=tmp/vischeck.coverage.out ./...
	go tool cover -func=tmp/vischeck.coverage.out

test-html: test
	go tool cover -html=tmp/vischeck.coverage.out -o tmp/vischeck.coverage.html
	open tmp/vischeck.coverage.html
