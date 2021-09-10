PKG = src

.PHONY: init

init:
	go mod download

format:
	gofmt ${PKG}/

run:
	go run main.go

build:
	go build .