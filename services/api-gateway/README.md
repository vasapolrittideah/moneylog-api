# api gateway

The API Gateway serves as the single entry point for all client requests in the moneylog-api microservices architecture. It handles request routing, authentication, rate limiting, and provides a unified HTTP REST API interface while communicating with backend services via gRPC.

## Architecture

The API Gateway has the following structure:

```
services/api-gateway/
├── cmd/                   # Application entry points
│   └── main.go            # Main application bootstrapper
├── internal/              # Private application code
│   ├── config/            # Application configuration
│   ├── handler/           # HTTP request handlers and routing
│   ├── middleware/        # HTTP middleware (auth, logging, CORS, etc.)
│   ├── payload/           # Request/response payload definitions
│   └── validator/         # Request validation logic
├── docs/                  # API documentation and OpenAPI specs
└── README.md              # Service documentation
```

### Responsibilities

1. **Handler** (`internal/handler/`)
   - HTTP request routing and handling
   - Request/response transformation between HTTP and gRPC
   - API versioning and endpoint management
   - Error handling and response formatting

2. **Middleware** (`internal/middleware/`)
   - Authentication and authorization
   - Request logging and metrics collection
   - Rate limiting and throttling
   - CORS handling and security headers
   - Request/response validation

3. **Payload** (`internal/payload/`)
   - HTTP request/response data structures
   - JSON serialization/deserialization
   - API contract definitions
   - Data transformation utilities

4. **Validator** (`internal/validator/`)
   - Input validation rules and logic
   - Business rule validation
   - Schema validation for API requests
   - Custom validation functions

## Key Features

- **Unified API Interface**: Single HTTP REST endpoint for all client interactions
- **Service Discovery**: Automatic discovery and load balancing via Consul
- **Protocol Translation**: HTTP to gRPC communication with backend services
- **Authentication**: JWT-based authentication and authorization
- **Rate Limiting**: Configurable rate limiting per client/endpoint
- **Observability**: Structured logging, metrics, and request tracing
