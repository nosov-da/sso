# SSO Service

A Single Sign-On (SSO) service built with Domain-Driven Design (DDD) architecture in Go.

## Architecture

This project follows Domain-Driven Design principles with clear separation of concerns:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: Repositories, event publishers, and external services
- **API Layer**: HTTP handlers and endpoints

See [DDD_ARCHITECTURE.md](DDD_ARCHITECTURE.md) for detailed architecture documentation.

## Prerequisites

- Go 1.24.4 or higher
- Make (optional, for convenience commands)
- Docker (for containerized deployment)

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd sso
```

### 2. Setup environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Build the application

```bash
make build
# or
go build -o bin/sso ./cmd/sso
```

### 5. Run the application

```bash
make run
# or
go run ./cmd/sso
```

The server will start on the address specified in your `.env` file (default: `:8080`).

## API Endpoints

### Register User

```bash
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response (201 Created):
```json
{
  "message": "user registered successfully"
}
```

### Login

```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Health Check

```bash
GET /health
```

Response (200 OK):
```json
{
  "status": "ok"
}
```

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SSO_ADDR` | Server address and port | `:8080` |
| `SSO_JWT_SECRET` | JWT signing secret (required) | - |
| `SSO_TOKEN_TTL` | Token TTL in seconds | `3600` |
| `SSO_DB_DSN` | Database connection string | `postgres://...` |
| `SSO_LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `ENVIRONMENT` | Environment name | `development` |

## Development

### Running Tests

```bash
make test
```

### Coverage Report

```bash
make coverage
```

Generates an HTML coverage report in `coverage.html`.

### Code Formatting

```bash
make fmt
```

### Linting

Install golangci-lint first:
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

Then run:
```bash
make lint
```

## Docker

Build and run with Docker:

```bash
make docker-build
make docker-run
```

Or manually:

```bash
docker build -t sso:latest .
docker run -p 8080:8080 \
  -e SSO_JWT_SECRET="your-secret-key" \
  sso:latest
```

## Project Structure

```
sso/
├── cmd/
│   └── sso/
│       └── main.go              # Application entry point
├── domain/
│   └── user/
│       ├── entity.go            # User entity
│       ├── value_objects.go     # Email, Password, UserId, Token
│       ├── repository.go        # Repository interface
│       └── events.go            # Domain events
├── application/
│   └── user/
│       └── use_cases.go         # Use cases (business logic)
├── infrastructure/
│   ├── persistence/             # Repository implementations
│   ├── events/                  # Event publisher implementations
│   └── token/                   # Token service implementations
├── internal/
│   ├── api/                     # HTTP handlers
│   ├── auth/                    # Authentication logic
│   ├── config/                  # Configuration
│   ├── middleware/              # HTTP middleware
│   └── server/                  # HTTP server setup
└── migrations/                  # Database migrations
```

## Best Practices

### Error Handling

The codebase uses domain-specific error types defined in `domain/user/errors.go`. Always use these predefined errors for consistency and better error handling.

### Logging

The application uses structured logging with `log/slog`. All logs are in JSON format for easy parsing by log aggregation tools.

### Password Security

Passwords are hashed using bcrypt with default cost factor. Never store plain text passwords.

### Validation

Input validation is performed at both application and domain layers:
- Application layer: Input presence and format
- Domain layer: Business rules enforcement

## Contributing

1. Create a feature branch
2. Make your changes
3. Run `make fmt` and `make test`
4. Commit with clear messages
5. Create a pull request

## License

See LICENSE file for details.
