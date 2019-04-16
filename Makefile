.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o build/router

.PHONY: test
test:
	go test ./...
