
.PHONY: fmt build


build:
	@go build -o ./bin/warden ./cmd/warden/
	@chmod +x ./bin/warden

fmt:
	@go mod tidy -v
	@go fmt ./...
