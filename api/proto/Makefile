# Proto Files Build Configuration
# Usage: make proto

PROTO_DIR := .
PROTO_LIB_DIR := ../../proto-lib
PROTO_MODULE := github.com/rusgainew/tunduck-app-mk/proto-lib

PROTOC_BASE := \
	protoc --proto_path=$(PROTO_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

# Legacy proto files (organized in legacy/ directory)
PROTO_LEGACY_AUTH := $(wildcard legacy/auth*.proto)
PROTO_LEGACY_COMPANY := $(wildcard legacy/company*.proto)
PROTO_LEGACY_DOCUMENT := $(wildcard legacy/document*.proto)
PROTO_LEGACY_COMMON := $(wildcard legacy/common.proto)

# New organized proto files
PROTO_DICTIONARIES := $(wildcard dictionaries/*.proto)
PROTO_ENTITIES := $(wildcard entities/*.proto)
PROTO_API := $(wildcard api/*.proto)
PROTO_ERRORS := $(wildcard errors/*.proto)

.PHONY: proto
proto: proto-clean proto-dirs proto-legacy proto-new

.PHONY: proto-legacy
proto-legacy:
	@echo "Building legacy proto files..."
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) legacy/common.proto
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) legacy/auth.proto legacy/auth_service.proto
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) legacy/company.proto legacy/company_service.proto
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) legacy/document.proto legacy/document_service.proto

.PHONY: proto-new
proto-new:
	@echo "Building new organized proto files..."
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_DICTIONARIES)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_ENTITIES)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_API)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_ERRORS)

.PHONY: proto-dirs
proto-dirs:
	mkdir -p $(PROTO_LIB_DIR)/common $(PROTO_LIB_DIR)/auth $(PROTO_LIB_DIR)/company $(PROTO_LIB_DIR)/document
	mkdir -p $(PROTO_LIB_DIR)/dictionaries $(PROTO_LIB_DIR)/entities $(PROTO_LIB_DIR)/api $(PROTO_LIB_DIR)/errors

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
