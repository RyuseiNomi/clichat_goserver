GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=clichat
BINARY_UNIX=$(BINARY_NAME)_unix

build:
	$(GOBUILD) -o $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

exec-dev:
	$(GOBUILD) -o $(BINARY_NAME); ./$(BINARY_NAME)