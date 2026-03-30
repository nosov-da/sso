# SSO Project Structure

## Complete Project Layout (with DDD)

```
sso/
├── domain/                           # Domain Layer (Pure Business Logic)
│   └── user/
│       ├── value_objects.go         # Email, Password, UserId, Token value objects
│       ├── value_objects_test.go    # Value object tests
│       ├── entity.go                # User entity
│       ├── entity_test.go           # Entity tests
│       ├── repository.go            # Repository interface (abstraction)
│       ├── events.go                # Domain events
│       ├── events_test.go           # Event tests
│
├── application/                      # Application Layer (Use Cases)
│   └── user/
│       ├── use_cases.go             # RegisterUserUseCase, AuthenticateUserUseCase
│       └── use_cases_test.go        # Use case tests
│
├── infrastructure/                   # Infrastructure Layer (Technical Implementation)
│   ├── persistence/
│   │   ├── in_memory_user_repository.go    # In-memory repository implementation
│   │   └── in_memory_user_repository_test.go
│   ├── events/
│   │   ├── event_bus.go             # EventPublisher implementations
│   │   └── event_bus_test.go
│   └── token/
│       └── jwt_service.go           # JWT token service
│
├── internal/                         # Internal packages (Legacy + Adapters)
│   ├── api/
│   │   ├── handler.go               # HTTP handlers (Adapter layer)
│   │   └── handler_test.go          # Handler tests
│   ├── auth/
│   │   ├── jwt.go                   # JWT manager implementation
│   │   └── jwt_test.go              # JWT tests
│   ├── config/
│   │   └── config.go                # Configuration loading
│   ├── server/
│   │   └── server.go                # HTTP server setup
│   ├── model/
│   │   └── user.go                  # Legacy user model
│   ├── repository/
│   │   ├── user_repo.go             # Legacy repository interface
│   │   ├── memory_user_repo.go      # Legacy in-memory implementation
│   │   └── memory_user_repo_test.go
│   └── service/
│       ├── auth_service.go          # Legacy auth service
│       └── auth_service_test.go
│
├── cmd/
│   └── sso/
│       └── main.go                  # Application entry point
│
├── migrations/
│   └── 0001_create_users.sql        # Database migrations
│
├── docs/
│   └── openapi.yaml                 # API documentation
│
├── go.mod
├── go.sum
├── docker-compose.yml
├── Dockerfile
├── DDD_ARCHITECTURE.md               # DDD architecture documentation
└── PROJECT_STRUCTURE.md              # This file
```

## Layer Responsibilities

### Domain Layer (`domain/user/`)
**Purpose:** Pure business logic, independent of any framework or technology

- **value_objects.go**
  - `Email`: Validated email address (cannot be empty, must be valid format)
  - `Password`: Hashed password (never stores plain text)
  - `UserId`: User identifier
  - `Token`: Authentication token

- **entity.go**
  - `User`: Domain entity representing a user
  - Methods: `ID()`, `Email()`, `Password()`, `Equals()`, `IdentifyBy()`

- **repository.go**
  - `Repository`: Interface defining persistence contracts
  - Methods: `Save()`, `GetByEmail()`, `GetByID()`, `Exists()`

- **events.go**
  - `UserRegisteredEvent`: Fired when user registers
  - `UserLoggedInEvent`: Fired when user logs in
  - `EventPublisher`: Interface for publishing events

### Application Layer (`application/user/`)
**Purpose:** Orchestrate domain logic, implement use cases

- **use_cases.go**
  - `RegisterUserUseCase`: Register new user
    - Input: `RegisterInput{Email, Password}`
    - Process: Validate, hash, save, publish event
    - Output: Error or success
    
  - `AuthenticateUserUseCase`: Authenticate user
    - Input: `AuthenticateInput{Email, Password}`
    - Process: Validate, verify password, generate token, publish event
    - Output: `AuthenticateOutput{Token}` or error

### Infrastructure Layer (`infrastructure/`)
**Purpose:** Technical implementations of domain contracts

- **persistence/in_memory_user_repository.go**
  - Implements `domain.user.Repository`
  - Thread-safe in-memory storage with RWMutex

- **events/event_bus.go**
  - `InMemoryEventPublisher`: Subscribes handlers, publishes events in-memory
  - `NullEventPublisher`: No-op publisher for testing

- **token/jwt_service.go**
  - `JWTTokenService`: Wraps JWT manager to implement TokenService interface

### Adapter Layer (`internal/api/`)
**Purpose:** Converts HTTP protocol to domain operations

- **handler.go**
  - `AuthHandler`: HTTP handlers for `/register` and `/login`
  - Receives HTTP requests, converts to use case inputs
  - Calls use cases, converts outputs to HTTP responses

### Support/Legacy Layer (`internal/`)
**Purpose:** Shared utilities and legacy code

- **auth/jwt.go**: JWT token management
- **config/config.go**: Environment-based configuration
- **server/server.go**: HTTP server setup with graceful shutdown
- **repository/**: Legacy repository implementations (kept for compatibility)
- **service/**: Legacy auth service (kept for compatibility)

## Data Flow Example: User Registration

```
HTTP Request (POST /register)
    ↓
[Handler.Register()] - Parse JSON
    ↓
RegisterUserUseCase.Execute(RegisterInput)
    ├─ NewEmail(email) - Value Object creation with validation
    ├─ repository.Exists(email) - Check if user exists
    ├─ bcrypt.GenerateFromPassword() - Hash password
    ├─ NewPassword(hash) - Value Object creation
    ├─ NewUserId(id) - Value Object creation
    ├─ User.NewUser(userId, email, password) - Create entity
    ├─ repository.Save(user) - Persist entity
    ├─ eventPublisher.Publish(UserRegisteredEvent) - Publish event
    └─ Return nil or error
    ↓
[Handler] - Send HTTP response (201 or 400)
    ↓
HTTP Response
```

## Testing Strategy

### Domain Layer Tests
- Test value object validation
- Test entity identity and equality
- Test domain invariants

### Application Layer Tests
- Test use case execution
- Test business logic coordination
- Test event publishing

### Infrastructure Layer Tests
- Test repository operations
- Test event publishing
- Test external service implementations

### Adapter Layer Tests
- Test HTTP request/response handling
- Test integration with use cases
- Test error handling and status codes

## Key Design Patterns Used

1. **Value Objects**: Immutable objects that represent specific concepts
2. **Entities**: Objects with unique identity and persistence
3. **Repository Pattern**: Abstraction for data access
4. **Domain Events**: Publish events that represent domain state changes
5. **Application Services**: Orchestrate domain logic for use cases
6. **Dependency Injection**: Inject dependencies through constructors
7. **Interface Segregation**: Small, focused interfaces for flexibility
8. **Hexagonal Architecture**: Isolated domain with adapters for external communication

## Moving Between Layers

### From External (HTTP) -> To Domain
```
HTTP Request → Handler → Application Service → Domain Entity
```

### From Domain -> To External (Repository)
```
Domain Entity → Repository Interface → Infrastructure Implementation
```

### Event Publishing
```
Domain Event → EventPublisher Interface → Infrastructure Handler
```

## Adding New Features

To add a new feature (e.g., password reset):

1. **Domain Layer**: Create value objects and entities
2. **Domain Events**: Define `UserPasswordResetEvent`
3. **Repository**: Add methods like `GetByEmail()`
4. **Application**: Create `ResetPasswordUseCase`
5. **Infrastructure**: Implement repository methods, event handlers
6. **Adapter**: Add HTTP handlers for new endpoints
7. **Tests**: Add tests at each layer

## Legacy Code Compatibility

The project maintains partial legacy code:
- `internal/service/auth_service.go`: Old service implementation
- `internal/repository/`: Old repository implementations

These are kept for reference but the new code uses DDD architecture exclusively. You can gradually migrate or remove legacy code as needed.

## Configuration

Set environment variables:
```bash
SSO_ADDR=:8080                          # Server address
SSO_JWT_SECRET=your-secret-key          # Required: JWT signing secret
SSO_DB_DSN=postgres://...               # Optional: PostgreSQL connection
```

## Building and Running

```bash
# Build
go build -o sso cmd/sso/main.go

# Run
export SSO_JWT_SECRET="your-secret"
./sso

# Test
go test ./...

# Test with coverage
go test ./... -cover
```
