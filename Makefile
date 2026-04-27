GOCACHE ?= /tmp/go-build-cache
BINARY_DIR := bin
BINARY := $(BINARY_DIR)/habit-examiner

.PHONY: build run-static run-full run-no-e2e cross

build:
	mkdir -p $(BINARY_DIR)
	GOCACHE=$(GOCACHE) go build -o $(BINARY) ./examiner

run-static: build
	./$(BINARY) -repo . -run-commands=false

run-full: build
	./$(BINARY) -repo .

run-no-e2e: build
	./$(BINARY) -repo . -run-e2e=false

cross:
	mkdir -p $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 GOCACHE=$(GOCACHE) go build -o $(BINARY_DIR)/habit-examiner-linux-amd64 ./examiner
	GOOS=darwin GOARCH=arm64 GOCACHE=$(GOCACHE) go build -o $(BINARY_DIR)/habit-examiner-darwin-arm64 ./examiner
	GOOS=windows GOARCH=amd64 GOCACHE=$(GOCACHE) go build -o $(BINARY_DIR)/habit-examiner-windows-amd64.exe ./examiner
