.PHONY: proto

# This command generates the Go code from the .proto files.
proto:
	@echo "Generating Go code from protobuf..."
	@protoc --proto_path=internal/proto \
	       --go_out=gen/ --go_opt=paths=source_relative \
	       --go-grpc_out=gen/ --go-grpc_opt=paths=source_relative \
	       internal/proto/token/token.proto
	@echo "Done."

# You can add more commands here, like 'build', 'test', 'run', etc.
build:
	@echo "Building the application..."
	@go build -o ./bin/auth-service ./cmd/server

run: proto build
	@echo "Running the application..."
	@./bin/auth-service