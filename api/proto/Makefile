# Proto Files Build Configuration
# Usage: make proto

PROTO_DIR := .
PROTO_LIB_DIR := ../../proto-lib
PROTO_MODULE := github.com/rusgainew/tunduck-app-mk/proto-lib

PROTOC_BASE := \
	protoc --proto_path=$(PROTO_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

PROTO_AUTH := auth.proto auth_service.proto
PROTO_COMPANY := company.proto company_service.proto
PROTO_DOCUMENT := document.proto document_service.proto
PROTO_COMMON := common.proto

.PHONY: proto
proto: proto-clean proto-dirs
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR)/common --go-grpc_out=$(PROTO_LIB_DIR)/common $(PROTO_COMMON)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR)/auth --go-grpc_out=$(PROTO_LIB_DIR)/auth $(PROTO_AUTH)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR)/company --go-grpc_out=$(PROTO_LIB_DIR)/company $(PROTO_COMPANY)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR)/document --go-grpc_out=$(PROTO_LIB_DIR)/document $(PROTO_DOCUMENT)

.PHONY: proto-dirs
proto-dirs:
	mkdir -p $(PROTO_LIB_DIR)/common $(PROTO_LIB_DIR)/auth $(PROTO_LIB_DIR)/company $(PROTO_LIB_DIR)/document

.PHONY: proto-clean
proto-clean:
	@[ -d $(PROTO_LIB_DIR) ] && find $(PROTO_LIB_DIR) -name "*.pb.go" -delete || true
	find $(PROTO_DIR) -maxdepth 1 -name "*.pb.go" -delete

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
