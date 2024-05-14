.DEFAULT_GOAL := build

build: vet
	@go build -o frtc .
.PHONY: build

fmt:
	@go fmt ./...
.PHONY: fmt

lint: fmt
	@golint ./...
.PHONY: lint

vet: fmt
	@go vet ./...
.PHONY: vet

build-all: vet
	@echo "compiling windows version ..."
	@GOOS=windows GOARCH=amd64 go build -o ./bin/win/frtc.exe .

	@echo "compiling linux version ..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/linux/frtc .

	@echo "compiling macos version ..."
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/macos/frtc .
.PHONY: build-all

install: build
	@sudo mv ./frtc /bin/
.PHONY: install