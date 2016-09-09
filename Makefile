build:
        go build -o farmer main.go

deps:
        go get ./...

.PHONY: build deps
