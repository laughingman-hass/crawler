.PHONY: build
build:
	go build ./cmd/crawler/crawler.go

.PHONY: test
test:
	go test -v ./...
