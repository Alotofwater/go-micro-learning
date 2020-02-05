
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o user-srv main.go lib.go

.PHONY: docker
docker:docker
	docker build . -t user-srv:latest
