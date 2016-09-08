build:
	go build -o farmer main.go

deps:
	go get -u github.com/tools/godep
	godep restore

.PHONY: build deps
