GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=dht-prometheus-exporter
BINARY_DEST=/usr/bin
BUILD_DIR=./cmd/dht-prometheus-exporter

.PHONY: all
all: test build

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(BUILD_DIR)

.PHONY: test
test:
	$(GOTEST) -v -race -cover ./...

.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

.PHONY: install
install: build
	sudo cp -f $(BINARY_NAME) $(BINARY_DEST)

.PHONY: uninstall
uninstall:
	sudo rm -f $(BINARY_DEST)/$(BINARY_NAME)

.PHONY: mod-tidy
mod-tidy:
	$(GOCMD) mod tidy
