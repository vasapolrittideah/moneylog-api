.PHONY: create-service
create-service:
	@if [ -z "$(name)" ]; then \
		echo "Please provide a service name using name variable"; \
		echo "Usage: make create-service name=<service-name>"; \
		echo "Example: make create-service name=user"; \
		exit 1; \
	else \
		./scripts/create_service.sh -name $(name); \
	fi

PROTO_DIR := protos
PROTO_SRC := $(shell find $(PROTO_DIR) -name "*.proto")
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	@mkdir -p shared/protos
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)
