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
