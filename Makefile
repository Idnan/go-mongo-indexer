.PHONY: build mac linux windows all

default: build

build:
	go build -o bin/indexer cmd/indexer/*.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/indexer cmd/indexer/*.go

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/indexer cmd/indexer/*.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/indexer.exe cmd/indexer/*.go

all: mac linux windows