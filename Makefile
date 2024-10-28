
.PHONY: fmt build debug


build:
	@go build -o ./bin/warden ./cmd/warden/
	@chmod +x ./bin/warden

debug:
	@go build -gcflags=all="-N -l" -o ./bin/warden ./cmd/warden
	@chmod +x ./bin/warden
	@./bin/warden init -p ./tmp/test

fmt:
	@go mod tidy -v
	@go fmt ./...

lint:
	@golangci-lint run ./...
