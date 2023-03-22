protoc: protoc
	@echo "Generating Go files from proto"
	mkdir -p build & cd proto && protoc --go_out=plugins=grpc:../build -I. *.proto

server: protoc
	@echo "Building server"
	go build -o bin/server src/server.go

all: protoc server

clean: clean
	@echo "Clenaing..."
	rm -rf build && rm -rf bin

.PHONY: server protoc clean
