
.PHONY: fmt build


build:
	@go build -o bin/warden cmd/warden/main.go

fmt:
	@go mod tidy -v
	@go fmt ./...
