GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gateway-api main.go lib.go

.PHONY: docker
docker:
	docker build . -t auth-srv:latest
