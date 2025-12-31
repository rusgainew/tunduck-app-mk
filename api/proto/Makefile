# Proto Files Build Configuration
# Usage: make proto

.PHONY: proto
proto: proto-clean
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		api/proto/*.proto

.PHONY: proto-clean
proto-clean:
	find . -name "*.pb.go" -delete

.PHONY: proto-install-tools
proto-install-tools:
	@echo "Installing protobuf compiler and Go plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Proto tools installed successfully!"

.PHONY: proto-check
proto-check:
	@echo "Checking proto files..."
	protoc --version
	protoc-gen-go --version
	protoc-gen-go-grpc --version

.PHONY: help
help:
	@echo "Proto Files Management"
	@echo "====================="
	@echo ""
	@echo "Available targets:"
	@echo "  make proto              - Compile all .proto files"
	@echo "  make proto-clean        - Remove all generated .pb.go files"
	@echo "  make proto-install-tools - Install protobuf compiler and Go plugins"
	@echo "  make proto-check        - Check installed proto tools versions"
	@echo "  make help               - Show this help message"
	@echo ""
	@echo "Proto files location: api/proto/"
	@echo "Generated code: *.pb.go and *_grpc.pb.go in the same directory"
