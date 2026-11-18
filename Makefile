all: protoc build

build: cmd/convoy/main.go cmd/convoy-client/main.go pkg/client.go pkg/server.go pkg/db.go pkg/constants.go go.mod go.sum
	go build -C cmd/convoy -buildvcs=false; cd cmd/convoy; install
	go build -C cmd/convoy-client -buildvcs=false; cd cmd/convoy-client; install

protoc: proto/convoy.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/convoy.proto proto/worker.proto --descriptor_set_out=docs/descriptor.pb --include_source_info
	cp proto/worker*.pb.go ../gown/proto/worker
	cp proto/convoy*.pb.go ../gown/proto/convoy
