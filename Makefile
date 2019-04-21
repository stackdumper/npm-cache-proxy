.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	gox -output=build/ncp_{{.OS}}_{{.Arch}}

.PHONY: test
test:
	go test ./...
