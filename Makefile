# Makefile

APP_NAME := bin/main
GO_SRC := $(shell find . -name "*.go" -type f)

build: $(APP_NAME)

$(APP_NAME): $(GO_SRC)
	go build -o $(APP_NAME) main.go

clean:
	rm -rf bin/*

.PHONY: build clean
