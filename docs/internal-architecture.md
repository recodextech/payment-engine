# Internal Folder Architecture Documentation

## Overview

The `internal` folder follows a clean architecture pattern with clear separation of concerns, implementing a layered architecture that promotes maintainability, testability, and scalability. This architecture follows Domain-Driven Design (DDD) principles and Clean Architecture patterns.

## Architecture Layers

### 1. HTTP Layer (`/internal/http`)

**Purpose**: Handles HTTP request and response processing, routing, and web server configuration.

**Components**:
- **Request Handlers**: Process incoming HTTP requests
- **Request/Response Models**: Define data structures for HTTP communication
- **Routing**: Configure URL paths and HTTP methods
- **Middleware**: Handle cross-cutting concerns like logging, authentication
- **Error Handling**: Format and return appropriate HTTP error responses

**Key Files**:
- `server.go` - HTTP server configuration and startup
- `router.go` - Route definitions and middleware setup
- `routes.go` - Specific route handlers mapping
- `handlers/` - Individual request handlers
- `request/` - Request DTOs and validation
- `responses/` - Response models and writers
- `errors/` - HTTP error handling and formatting

**Flow**:
```
HTTP Request → Router → Handler → Service Layer → Response
```

### 2. Domain Layer (`/internal/domain`)

**Purpose**: Contains core business logic, domain models, interfaces, and business rules. This is the heart of the application that defines what the system does.

**Components**:
- **Domain Models**: Core business entities and value objects
- **Domain Events**: Events that represent business state changes
- **Domain Services**: Business logic that doesn't belong to a specific entity
- **Interfaces/Contracts**: Define contracts for external dependencies
- **Domain Definitions**: Constants, enums, and domain-specific types

**Key Directories**:
- `models/` - Domain entities and value objects
- `events/` - Domain events for state changes
- `service/` - Domain service interfaces
- `adaptors/` - Interface definitions for external dependencies
- `application/` - Application-level configurations and modules
- `definitions.go` - Domain constants and enums
- `tables.go` - Database table name definitions

**Domain Events Examples**:
- `ProcessCreated` - When a new process is created
- `JobCreated` - When a job is assigned to a process
- `JobLocationChanged` - When job locations are updated

### 3. Adaptors Layer (`/internal/adaptors`)

**Purpose**: Handles external communications and provides implementations for domain interfaces. Acts as the bridge between the domain and external systems.

**Components**:
- **Database Adaptors**: Database connection and query implementations
- **Repository Implementations**: Data access layer implementations
- **Encoders**: Data serialization/deserialization (JSON, Avro)
- **External Service Clients**: Third-party service integrations

**Key Directories**:
- `database/postgres/` - PostgreSQL database implementations
- `repositories/` - Repository pattern implementations
- `encoders/` - Data encoding/decoding implementations
  - `json/` - JSON encoding
  - `avro/` - Avro schema registry integration
  - `misc/` - Utility encoders (UUID, string)

**Repository Pattern Implementation**:
- `ProcessRepository` - Manages process-related data operations
- `JobRepository` - Handles job data persistence
- `JobLocationRepository` - Manages job location data

### 4. App Layer (`/internal/app`)

**Purpose**: Contains application-level configurations, settings, and global application state management.

**Components**:
- **Configuration Management**: Environment-based configurations
- **Application Settings**: Debug modes, feature flags
- **Global State**: Application-wide state management

**Key Files**:
- `config.go` - Application configuration and environment settings

### 5. Bootstrap Layer (`/internal/bootstrap`)

**Purpose**: Service startup orchestration, dependency injection, and application lifecycle management. This layer ensures proper initialization order and binds all components together.

**Components**:
- **Dependency Injection**: Wire up all dependencies
- **Service Registration**: Register services in the container
- **Initialization Order**: Ensure proper startup sequence
- **Lifecycle Management**: Handle application start and shutdown

**Key Files**:
- `boot.go` - Main bootstrap entry point with configuration setup
- `bind.go` - Dependency injection and service binding
- `init.go` - Module initialization orchestration
- `start.go` - Service startup and shutdown handling

**Bootstrap Process**:
1. **Configuration Binding** - Load and bind all configurations
2. **Service Binding** - Register all services in the DI container
3. **Module Initialization** - Initialize modules in proper order
4. **Service Startup** - Start HTTP server and background services

### 6. Services Layer (`/internal/services`)

**Purpose**: Orchestrates business operations by coordinating between HTTP handlers and use case layers. Handles cross-cutting concerns and business workflows.

**Components**:
- **Business Orchestration**: Coordinate multiple domain operations
- **Transaction Management**: Handle data consistency across operations
- **Business Validation**: Enforce business rules
- **Service Coordination**: Manage interactions between different domains

**Key Files**:
- `trip_create_service.go` - Trip creation business logic
- `errors.go` - Service-level error definitions

**Service Responsibilities**:
- Validate business rules
- Coordinate multiple repository operations
- Handle transactions
- Transform between HTTP models and domain models

### 7. Use Cases Layer (`/internal/usecases`)

**Purpose**: Implements specific application use cases and business scenarios. Encapsulates application-specific business logic that orchestrates domain services and repositories.

**Components**:
- **Application Logic**: Use case specific implementations
- **Workflow Orchestration**: Manage complex business processes
- **Data Transformation**: Convert between layers
- **Business Process Management**: Handle multi-step operations

**Key Files**:
- `trip_create_usecase.go` - Trip creation use case implementation

## Data Flow Architecture

### Request Flow
```
HTTP Request
    ↓
HTTP Router
    ↓
Request Handler
    ↓
Service Layer
    ↓
Use Case Layer
    ↓
Domain Layer
    ↓
Repository (Adaptor)
    ↓
Database
```

### Response Flow
```
Database
    ↓
Repository (Adaptor)
    ↓
Domain Layer
    ↓
Use Case Layer
    ↓
Service Layer
    ↓
Response Handler
    ↓
HTTP Response
```

## Layer Dependencies

### Dependency Rule
- **Outer layers depend on inner layers**
- **Inner layers never depend on outer layers**
- **Domain layer has no external dependencies**

### Dependency Graph
```
HTTP Layer
    ↓
Services Layer
    ↓
Use Cases Layer
    ↓
Domain Layer
    ↑
Adaptors Layer
```

## Key Design Patterns

### 1. Repository Pattern
- Abstracts data access logic
- Provides a uniform interface for data operations
- Enables testing with mock implementations

### 2. Dependency Injection
- Managed through `recodextech/container`
- Enables loose coupling between components
- Facilitates testing and maintainability

### 3. Domain Events
- Represents business state changes
- Enables event-driven architecture
- Supports audit trails and integration patterns

### 4. Clean Architecture
- Separation of concerns across layers
- Business logic isolated from infrastructure
- Testable and maintainable codebase

## Configuration Management

### Environment-Based Configuration
- Each layer has its own configuration structure
- Environment variables control behavior
- Configuration validation at startup

### Configuration Layers:
- **HTTP Server**: Port, timeouts, routing settings
- **Database**: Connection parameters, pool settings
- **Logging**: Log levels, output formats
- **Metrics**: Prometheus metrics configuration
- **Schema Registry**: Avro schema management

## Error Handling Strategy

### Layer-Specific Error Handling:
1. **HTTP Layer**: HTTP status codes and error responses
2. **Service Layer**: Business error codes and messages
3. **Domain Layer**: Domain-specific errors and validations
4. **Adaptor Layer**: Infrastructure and external service errors

### Error Propagation:
- Errors bubble up through layers
- Each layer adds context and transforms errors appropriately
- Consistent error format across the application

## Testing Strategy

### Layer Testing:
- **Unit Tests**: Domain logic and business rules
- **Integration Tests**: Database and external service interactions
- **HTTP Tests**: API endpoint testing
- **Contract Tests**: Interface compliance testing

### Mock Strategy:
- Repository interfaces enable easy mocking
- Dependency injection supports test doubles
- Domain services can be tested in isolation

## Scalability Considerations

### Horizontal Scaling:
- Stateless HTTP handlers
- Database connection pooling
- Event-driven communication patterns

### Performance Optimization:
- Repository pattern enables caching strategies
- Connection pooling for database operations
- Metrics and monitoring integration

## Development Guidelines

### Adding New Features:
1. Define domain models and events
2. Create repository interfaces in domain layer
3. Implement repositories in adaptors layer
4. Create use cases for business logic
5. Implement service layer coordination
6. Add HTTP handlers and routes
7. Update bootstrap configuration

### Best Practices:
- Keep domain layer pure (no external dependencies)
- Use interfaces for all external dependencies
- Implement proper error handling at each layer
- Follow consistent naming conventions
- Document business rules in domain layer
- Use dependency injection for all components

## Module Dependencies

### External Dependencies:
- **Container**: Dependency injection framework
- **KRouter**: HTTP routing and middleware
- **PostgreSQL**: Database operations
- **Prometheus**: Metrics collection
- **Avro**: Schema registry integration

### Internal Module Communication:
- All communication through well-defined interfaces
- Event-driven patterns for loose coupling
- Configuration-based feature toggles

This architecture promotes maintainability, testability, and scalability while following established software engineering principles and patterns.