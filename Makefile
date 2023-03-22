protoc: protoc
	@echo "Generating Go files from proto"
	cd proto && protoc --go_out=plugins=grpc:. -I. *.proto

server: protoc
	@echo "Building server"
	go build -o server src/server.go

all: protoc server

.PHONY: server protoc
