# Domain-Driven Design (DDD) Architecture

## Overview

The SSO project has been refactored to follow Domain-Driven Design principles, providing a clean separation between business logic (domain) and technical implementation (infrastructure).

## Architecture Structure

```
domain/
├── user/
│   ├── value_objects.go       # Value Objects: Email, Password, UserId, Token
│   ├── entity.go              # User Entity with domain methods
│   ├── repository.go          # Repository interface (domain contract)
│   ├── events.go              # Domain Events and EventPublisher interface
│   └── *_test.go              # Domain layer tests

application/
├── user/
│   ├── use_cases.go           # Application Services (Use Cases)
│   └── *_test.go              # Application layer tests

infrastructure/
├── persistence/
│   ├── in_memory_user_repository.go    # Repository implementation
│   └── *_test.go
├── events/
│   ├── event_bus.go           # Event publisher implementations
│   └── *_test.go
└── token/
    └── jwt_service.go         # Token service implementation

internal/
├── api/
│   ├── handler.go             # HTTP handlers (Adapter layer)
│   └── *_test.go
├── server/
│   └── server.go              # HTTP server setup
├── config/
│   └── config.go              # Configuration
└── ...
```

## Key DDD Concepts

### 1. Value Objects
Immutable objects defined by their attributes, not their identity.

**Files:** `domain/user/value_objects.go`

- **Email**: Validated email address
- **Password**: Hashed password (stored as hash only)
- **UserId**: User identifier
- **Token**: Authentication token

```go
email, err := user.NewEmail("user@example.com")
// Email is validated on creation
```

### 2. Entities
Objects with a unique identity that persists over time.

**File:** `domain/user/entity.go`

- **User Entity**: Represents a user with ID, email, and password
  - Methods: `ID()`, `Email()`, `Password()`, `Equals()`, `IdentifyBy()`

```go
user := user.NewUser(userId, email, password)
```

### 3. Repository Pattern
Interface defining persistence contracts (Repository Pattern).

**File:** `domain/user/repository.go`

```go
type Repository interface {
    Save(user *User) error
    GetByEmail(email Email) (*User, error)
    GetByID(id UserId) (*User, error)
    Exists(email Email) (bool, error)
}
```

**Implementations:**
- `infrastructure/persistence/in_memory_user_repository.go` - In-memory storage

### 4. Domain Events
Events that capture something important that happened in the domain.

**File:** `domain/user/events.go`

- **UserRegisteredEvent**: Fired when a user registers
- **UserLoggedInEvent**: Fired when a user logs in
- **EventPublisher Interface**: Contracts for publishing events

```go
event := user.NewUserRegisteredEvent(userId, email)
publisher.Publish(event)
```

**Implementations:**
- `infrastructure/events/event_bus.go` - In-memory event bus with subscribers

### 5. Application Services (Use Cases)
Orchestrate domain logic and coordinate between domain and infrastructure.

**File:** `application/user/use_cases.go`

- **RegisterUserUseCase**: Handles user registration
  - Input: `RegisterInput{Email, Password}`
  - Validates input, creates user, publishes event
  
- **AuthenticateUserUseCase**: Handles user authentication
  - Input: `AuthenticateInput{Email, Password}`
  - Retrieves user, verifies password, generates token, publishes event
  - Output: `AuthenticateOutput{Token}`

```go
registerUC := NewRegisterUserUseCase(repository, publisher)
err := registerUC.Execute(RegisterInput{Email: "user@example.com", Password: "pwd"})

authenticateUC := NewAuthenticateUserUseCase(repository, tokenService, publisher)
output, err := authenticateUC.Execute(AuthenticateInput{...})
```

### 6. Infrastructure Layer
Technical implementations of domain contracts.

**Persistence:**
- `infrastructure/persistence/in_memory_user_repository.go`: Repository implementation

**Events:**
- `infrastructure/events/event_bus.go`: 
  - `InMemoryEventPublisher`: Publishes events to in-memory subscribers
  - `NullEventPublisher`: No-op publisher for testing

**Token Service:**
- `infrastructure/token/jwt_service.go`: Wraps JWT manager

### 7. Adapter Layer (HTTP)
Converts HTTP requests to use cases and use cases to HTTP responses.

**File:** `internal/api/handler.go`

- Uses application services (use cases) instead of domain services
- Handles HTTP protocol details
- Converts JSON to domain inputs and domain outputs to JSON

## Data Flow

```
HTTP Request
    ↓
[Handler] - Parse JSON, create application input
    ↓
[UseCase] - Execute business logic
    ├─ Validate domain objects (Value Objects)
    ├─ Retrieve domain entities (Repository)
    ├─ Perform domain operations
    └─ Publish domain events (EventPublisher)
    ↓
Domain Entities
    ├─ User (with Email, Password, UserId)
    └─ Events (UserRegistered, UserLoggedIn)
    ↓
[Handler] - Convert output to JSON
    ↓
HTTP Response
```

## Benefits of DDD Architecture

1. **Separation of Concerns**: Business logic is isolated from technical details
2. **Testability**: Domain, application, and infrastructure layers can be tested independently
3. **Flexibility**: Easy to swap implementations (in-memory → PostgreSQL)
4. **Maintainability**: Clear structure makes code easier to understand and modify
5. **Domain Focus**: Code directly reflects business concepts and language
6. **Event-Driven**: Domain events enable reactive patterns for future features
7. **Dependency Inversion**: Domain doesn't depend on infrastructure, only on abstractions

## Testing Strategy

### Domain Layer Tests
Test value objects and entities in isolation:
- Value object creation and validation
- Entity behavior and identity
- Invariants and rules

**Files:** `domain/user/*_test.go`

### Application Layer Tests
Test use cases with mock infrastructure:
- Use case execution
- Business logic coordination
- Event publishing verification

**Files:** `application/user/*_test.go`

### Infrastructure Layer Tests
Test implementations with minimal dependencies:
- Repository operations
- Event bus publishing
- Token service functionality

**Files:** `infrastructure/*/\*_test.go`

### Adapter Layer Tests (HTTP)
Test HTTP handlers with mock domain implementations:
- Request parsing
- Response formatting
- Status codes

**Files:** `internal/api/*_test.go`

## Configuration and Wiring

Main application entry point demonstrates dependency injection:

**File:** `cmd/sso/main.go`

```go
// Infrastructure
userRepository := persistence.NewInMemoryUserRepository()
eventPublisher := events.NewInMemoryEventPublisher()
tokenService := token.NewJWTTokenService(jwtManager)

// Application (Use Cases)
registerUseCase := appsvc.NewRegisterUserUseCase(repository, publisher)
authenticateUseCase := appsvc.NewAuthenticateUserUseCase(repository, tokenService, publisher)

// Adapter (HTTP)
handler := api.NewAuthHandler(registerUseCase, authenticateUseCase)
```

## Future Extensions

With this DDD architecture, it's easy to add:

1. **Event Handlers**: Subscribe to domain events for side effects
   - Send email confirmations
   - Log audit trails
   - Update analytics

2. **PostgreSQL Repository**: Implement `domain/user/repository.go` for PostgreSQL

3. **More Domain Services**: Add password recovery, profile management, etc.

4. **Query Objects**: Separate read models from write models (CQRS pattern)

5. **Domain Services**: Add password validation, email uniqueness rules, etc.

## Key Files Summary

| Layer | File | Purpose |
|-------|------|---------|
| Domain | `domain/user/value_objects.go` | Email, Password, UserId, Token value objects |
| Domain | `domain/user/entity.go` | User entity |
| Domain | `domain/user/repository.go` | Repository interface |
| Domain | `domain/user/events.go` | Domain events |
| Application | `application/user/use_cases.go` | Register and Authenticate use cases |
| Infrastructure | `infrastructure/persistence/in_memory_user_repository.go` | Repository implementation |
| Infrastructure | `infrastructure/events/event_bus.go` | Event publishers |
| Infrastructure | `infrastructure/token/jwt_service.go` | Token service |
| Adapter | `internal/api/handler.go` | HTTP handlers |
| Entry Point | `cmd/sso/main.go` | Application bootstrap |
