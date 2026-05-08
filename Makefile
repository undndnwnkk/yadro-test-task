SERVER_NAME=dns-server
CLIENT_NAME=dns-client
BUILD_DIR=bin

.PHONY: all build test clean run-server help

all: build

build:
	@echo "Building binary files"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(SERVER_NAME) cmd/server/main.go
	go build -o $(CLIENT_NAME) cmd/client/main.go

test:
	@echo "Running tests"
	go test -v ./internal/dns/...

run-server:
	@echo "Starting server with sudo (password required)"
	sudo ./$(BUILD_DIR)/$(SERVER_NAME)

clean:
	@echo "Cleaning binary files"
	rm -rf $(BUILD_DIR)
	rm -rf $(CLIENT_NAME)

help:
	@echo "Commands for production mode (/etc/resolv.conf):"
	@echo " 	make build 	- Build binary files"
	@echo " 	make run-server  - Run server as sudo"
	@echo " 	make test  - Run unit tests"
	@echo " 	make clean  - Remove binary files"