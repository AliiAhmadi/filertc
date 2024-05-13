.DEFAULT_GOAL := build

build: 
	@go build -o frtc .
.PHONY: build

build-all:
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