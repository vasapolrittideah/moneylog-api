# auth service

This service handles all auth-related operations in the system.

## Architecture

The service follows Clean Architecture principles with the following structure:

```
services/auth_service/
├── cmd/                   # Application entry points
│   └── main.go            # Main application bootstrapper
├── internal/              # Private application code
│   ├── config/            # Application configuration
│   ├── domain/            # Business domain models and interfaces
│   ├── handler/           # Request handlers
│   │   └── grpc/          # gRPC server handlers
│   ├── service/           # Business logic implementation
│   └── store/             # Data persistence layer
│       └── mongo/         # MongoDB data storage implementation
├── pkg/                   # Public API and shared code
│   ├── client/            # Client libraries
│   └── types/             # Shared type definitions
└── README.md              # Service documentation
```

### Layer Responsibilities

1. **Presentation Layer** (`internal/handler/`)
   - Handles incoming requests and responses
   - Translates between transport payloads and domain models
   - Input validation and error handling
   - Protocol-specific implementations (gRPC, HTTP, etc.)

2. **Domain Layer** (`internal/domain/`)
   - Contains business domain models and entities
   - Defines contracts for stores and services
   - Pure business logic with no external dependencies
   - Domain-specific errors and validation rules

3. **Application Layer** (`internal/service/`)
   - Implements business use cases and workflows
   - Orchestrates domain objects and store operations
   - Handles cross-cutting concerns (logging, metrics)
   - Enforces business rules and policies

4. **Data Layer** (`internal/store/`)
   - Handles data persistence and retrieval
   - Implements domain store interfaces
   - Manages database connections and transactions
   - Data mapping and query optimization

## Key Benefits

1. **Dependency Inversion**: Services depend on interfaces, not concrete implementations
2. **Separation of Concerns**: Each layer has a single, well-defined responsibility
3. **Testability**: Easy to mock dependencies and write comprehensive unit tests
4. **Maintainability**: Clear boundaries and contracts between components
5. **Flexibility**: Easy to swap implementations without affecting business logic
6. **Scalability**: Modular design supports independent scaling and deployment

### Configuration
The service uses environment variables for configuration. See `internal/config/` for available options.

