GOPATH:=$(shell go env GOPATH)

.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/micro/micro/v2/cmd/protoc-gen-micro

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/dmicrog.proto
	
.PHONY: build
build:
	go build -o dmicrog *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: alpine-build
alpine-build: proto
	CGO_ENABLED=0 go build -o dmicrog *.go

.PHONY: docker
docker: alpine-build
	docker build . -t dmicrog:v0.0.1

.PHONY: up
up: alpine-build
	docker-compose build
	docker-compose up -d

.PHONY: down
down:
	docker-compose down -v
