all: protoc build

build: main.go client.go server.go go.mod go.sum
	go build -buildvcs=false; install

protoc: proto/convoy.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/convoy.proto proto/worker.proto --descriptor_set_out=docs/descriptor.pb --include_source_info
	cp proto/worker*.pb.go ../gown/proto/worker
	cp proto/convoy*.pb.go ../gown/proto/convoy
