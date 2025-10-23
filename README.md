# Microservices Architecture

A modern microservices architecture built with Go, implementing Domain-Driven Design (DDD), CQRS pattern, and Dependency Injection using Wire. This project consists of Product Service and Basket Service with comprehensive management capabilities.

## Services Overview

### 🛍️ Product Service
Comprehensive product management with HTTP REST API and gRPC interfaces.

### 🛒 Basket Service  
Shopping basket management with Redis storage and gRPC product client integration.

## Architecture Overview

```mermaid
graph TB
    subgraph "External Layer"
        HTTP[HTTP API<br/>Port 8080]
        GRPC[gRPC API<br/>Port 50050]
        DB[(PostgreSQL<br/>Port 5432)]
    end
    
    subgraph "Interface Layer"
        HTTPHandler[HTTP Handlers]
        GRPCHandler[gRPC Handlers]
        Middleware[Middleware<br/>CORS, Logging, Metrics]
    end
    
    subgraph "Application Layer (CQRS)"
        CommandHandler[Command Handler]
        QueryHandler[Query Handler]
        UseCase[Use Cases]
        DTO[DTOs]
    end
    
    subgraph "Domain Layer"
        Entity[Product Entity]
        Repository[Repository Interface]
        DomainService[Domain Service]
    end
    
    subgraph "Infrastructure Layer"
        RepoImpl[Repository Implementation]
        Config[Configuration]
        Logger[Logging]
        Metrics[Prometheus Metrics]
    end
    
    HTTP --> HTTPHandler
    GRPC --> GRPCHandler
    HTTPHandler --> Middleware
    GRPCHandler --> Middleware
    
    HTTPHandler --> CommandHandler
    HTTPHandler --> QueryHandler
    GRPCHandler --> CommandHandler
    GRPCHandler --> QueryHandler
    
    CommandHandler --> UseCase
    QueryHandler --> UseCase
    UseCase --> DTO
    UseCase --> Repository
    
    Repository --> RepoImpl
    RepoImpl --> DB
    
    Config --> Logger
    Config --> Metrics
    
    classDef external fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef interface fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef application fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef domain fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef infrastructure fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    
    class HTTP,GRPC,DB external
    class HTTPHandler,GRPCHandler,Middleware interface
    class CommandHandler,QueryHandler,UseCase,DTO application
    class Entity,Repository,DomainService domain
    class RepoImpl,Config,Logger,Metrics infrastructure
```

## CQRS Pattern Implementation

```mermaid
graph LR
    subgraph "Commands (Write Operations)"
        CreateCmd[CreateProductCommand]
        UpdateCmd[UpdateProductCommand]
        DeleteCmd[DeleteProductCommand]
    end
    
    subgraph "Queries (Read Operations)"
        GetProduct[GetProductQuery]
        GetProducts[GetProductsQuery]
        GetTopExpensive[GetTopMostExpensiveQuery]
        GetLowStock[GetLowStockProductsQuery]
        GetByCategory[GetProductsByCategoryQuery]
    end
    
    subgraph "Handlers"
        CmdHandler[Command Handler]
        QueryHandler[Query Handler]
    end
    
    subgraph "Use Cases"
        ProductUseCase[Product Use Case]
    end
    
    subgraph "Repository"
        ProductRepo[Product Repository]
    end
    
    CreateCmd --> CmdHandler
    UpdateCmd --> CmdHandler
    DeleteCmd --> CmdHandler
    
    GetProduct --> QueryHandler
    GetProducts --> QueryHandler
    GetTopExpensive --> QueryHandler
    GetLowStock --> QueryHandler
    GetByCategory --> QueryHandler
    
    CmdHandler --> ProductUseCase
    QueryHandler --> ProductUseCase
    ProductUseCase --> ProductRepo
    
    classDef command fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef query fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef handler fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef usecase fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef repository fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    
    class CreateCmd,UpdateCmd,DeleteCmd command
    class GetProduct,GetProducts,GetTopExpensive,GetLowStock,GetByCategory query
    class CmdHandler,QueryHandler handler
    class ProductUseCase usecase
    class ProductRepo repository
```

## Technology Stack

```mermaid
graph TB
    subgraph "Backend"
        Go[Go 1.21]
        Gin[Gin Framework]
        GORM[GORM ORM]
        Wire[Wire DI]
    end
    
    subgraph "Database"
        PostgreSQL[PostgreSQL 15]
        Migrations[Auto Migrations]
    end
    
    subgraph "APIs"
        REST[REST API]
        GRPC[gRPC API]
        Proto[Protocol Buffers]
    end
    
    subgraph "Monitoring"
        Prometheus[Prometheus Metrics]
        Logrus[Structured Logging]
        Health[Health Checks]
    end
    
    subgraph "DevOps"
        Docker[Docker]
        Compose[Docker Compose]
        Scripts[Shell Scripts]
    end
    
    Go --> Gin
    Go --> GORM
    Go --> Wire
    GORM --> PostgreSQL
    PostgreSQL --> Migrations
    
    Gin --> REST
    Go --> GRPC
    GRPC --> Proto
    
    Go --> Prometheus
    Go --> Logrus
    Go --> Health
    
    Go --> Docker
    Docker --> Compose
    Go --> Scripts
    
    classDef backend fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef database fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef api fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef monitoring fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef devops fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    
    class Go,Gin,GORM,Wire backend
    class PostgreSQL,Migrations database
    class REST,GRPC,Proto api
    class Prometheus,Logrus,Health monitoring
    class Docker,Compose,Scripts devops
```

## Project Structure

```mermaid
graph TD
    subgraph "Root Directory"
        CMD[cmd/product/]
        INTERNAL[internal/product/]
        API[api/proto/]
        SCRIPTS[scripts/]
        DOCKER[dockerfiles/]
        MAKEFILE[Makefile]
    end
    
    subgraph "Internal Structure"
        DOMAIN[domain/]
        APPLICATION[application/]
        INFRASTRUCTURE[infrastructure/]
        INTERFACES[interfaces/]
    end
    
    subgraph "Domain Layer"
        ENTITY[entity/]
        REPO[repository/]
        SERVICE[service/]
    end
    
    subgraph "Application Layer"
        COMMAND[command/]
        QUERY[query/]
        HANDLER[handler/]
        USECASE[usecase/]
        DTO[dto/]
    end
    
    subgraph "Infrastructure Layer"
        PERSISTENCE[persistence/]
        CONFIG[config/]
        EXTERNAL[external/]
    end
    
    subgraph "Interface Layer"
        HTTP[http/]
        GRPC[grpc/]
    end
    
    CMD --> INTERNAL
    INTERNAL --> DOMAIN
    INTERNAL --> APPLICATION
    INTERNAL --> INFRASTRUCTURE
    INTERNAL --> INTERFACES
    
    DOMAIN --> ENTITY
    DOMAIN --> REPO
    DOMAIN --> SERVICE
    
    APPLICATION --> COMMAND
    APPLICATION --> QUERY
    APPLICATION --> HANDLER
    APPLICATION --> USECASE
    APPLICATION --> DTO
    
    INFRASTRUCTURE --> PERSISTENCE
    INFRASTRUCTURE --> CONFIG
    INFRASTRUCTURE --> EXTERNAL
    
    INTERFACES --> HTTP
    INTERFACES --> GRPC
    
    classDef root fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef internal fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef domain fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef application fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef infrastructure fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef interfaces fill:#e0f2f1,stroke:#00695c,stroke-width:2px
    
    class CMD,INTERNAL,API,SCRIPTS,DOCKER,MAKEFILE root
    class DOMAIN,APPLICATION,INFRASTRUCTURE,INTERFACES internal
    class ENTITY,REPO,SERVICE domain
    class COMMAND,QUERY,HANDLER,USECASE,DTO application
    class PERSISTENCE,CONFIG,EXTERNAL infrastructure
    class HTTP,GRPC interfaces
```

## API Endpoints

```mermaid
graph TB
    subgraph "HTTP REST API (Port 8080)"
        GET_PRODUCTS[GET /products]
        GET_PRODUCT[GET /products/:id]
        CREATE_PRODUCT[POST /products]
        UPDATE_PRODUCT[PUT /products/:id]
        DELETE_PRODUCT[DELETE /products/:id]
        GET_TOP5[GET /products/top-5]
        GET_TOP10[GET /products/top-10]
        GET_LOW_STOCK1[GET /products/low-stock-1]
        GET_LOW_STOCK10[GET /products/low-stock-10]
        GET_BY_CATEGORY[GET /products/category/:category]
        HEALTH[GET /health]
        METRICS[GET /metrics]
    end
    
    subgraph "gRPC API (Port 50050)"
        GRPC_GET[GetProduct]
        GRPC_CREATE[CreateProduct]
        GRPC_UPDATE[UpdateProduct]
        GRPC_DELETE[DeleteProduct]
        GRPC_LIST[ListProducts]
        GRPC_TOP[GetTopMostExpensiveProducts]
        GRPC_LOW[GetLowStockProducts]
        GRPC_CATEGORY[GetProductsByCategory]
    end
    
    classDef http fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef grpc fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    
    class GET_PRODUCTS,GET_PRODUCT,CREATE_PRODUCT,UPDATE_PRODUCT,DELETE_PRODUCT,GET_TOP5,GET_TOP10,GET_LOW_STOCK1,GET_LOW_STOCK10,GET_BY_CATEGORY,HEALTH,METRICS http
    class GRPC_GET,GRPC_CREATE,GRPC_UPDATE,GRPC_DELETE,GRPC_LIST,GRPC_TOP,GRPC_LOW,GRPC_CATEGORY grpc
```

## Development Workflow

```mermaid
graph LR
    subgraph "Development"
        DEV[make dev]
        BUILD[make build]
        TEST[make test]
        LINT[make lint]
    end
    
    subgraph "Database"
        MIGRATE[make db-migrate]
        SEED[make db-seed]
        BACKUP[make db-backup]
    end
    
    subgraph "Docker"
        DOCKER_BUILD[make docker-build]
        DOCKER_RUN[make docker-run]
        DOCKER_STOP[make docker-stop]
    end
    
    subgraph "Cleanup"
        CLEAN[make clean]
        CLEAN_ALL[make clean-all]
    end
    
    DEV --> BUILD
    BUILD --> TEST
    TEST --> LINT
    
    MIGRATE --> SEED
    SEED --> BACKUP
    
    DOCKER_BUILD --> DOCKER_RUN
    DOCKER_RUN --> DOCKER_STOP
    
    CLEAN --> CLEAN_ALL
    
    classDef dev fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef db fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef docker fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef cleanup fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    
    class DEV,BUILD,TEST,LINT dev
    class MIGRATE,SEED,BACKUP db
    class DOCKER_BUILD,DOCKER_RUN,DOCKER_STOP docker
    class CLEAN,CLEAN_ALL cleanup
```

## 🛒 Basket Service

### Features
- **Redis Storage**: Fast in-memory storage with automatic TTL (1 day expiration)
- **gRPC Product Client**: Communicates with Product Service for product information
- **CQRS Pattern**: Command and Query separation for better performance
- **Simple Monitoring**: HTTP endpoint monitoring with Prometheus metrics
- **Auto Cleanup**: Redis TTL handles basket expiration automatically

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/baskets/{user_id}` | Get user's basket |
| POST | `/baskets` | Create new basket |
| POST | `/baskets/{user_id}/items` | Add item to basket |
| PUT | `/baskets/{user_id}/items/{product_id}` | Update item quantity |
| DELETE | `/baskets/{user_id}/items/{product_id}` | Remove item from basket |
| DELETE | `/baskets/{user_id}/items` | Clear all items |
| DELETE | `/baskets/{user_id}` | Delete entire basket |
| GET | `/health` | Health check |

### Environment Variables

```bash
# Service Configuration
PORT=8081
ENVIRONMENT=development

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Product Service Configuration
PRODUCT_SERVICE_URL=localhost:50050

# Logging Configuration
LOG_LEVEL=debug
LOG_FORMAT=text
LOG_OUTPUT=console
```

### Running Basket Service

```bash
# Build and run
go build -o bin/basket-service cmd/basket/main.go
./bin/basket-service

# Or with Docker Compose (includes Redis)
docker-compose up basket-service
```

### Example Usage

```bash
# Create a basket
curl -X POST http://localhost:8081/baskets \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123"}'

# Add item to basket
curl -X POST http://localhost:8081/baskets/user123/items \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 2}'

# Get basket
curl http://localhost:8081/baskets/user123

# Clear basket
curl -X DELETE http://localhost:8081/baskets/user123/items
```

## 🏗️ Development Workflow

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (for Product Service)
- Redis (for Basket Service)

### Quick Start

```bash
# Clone and setup
git clone <repository>
cd obs-tools-usage

# Start all services
docker-compose up -d

# Or start individual services
docker-compose up product-service
docker-compose up basket-service
```

### Service Communication

```mermaid
graph LR
    Client[Client] --> Product[Product Service<br/>:8080]
    Client --> Basket[Basket Service<br/>:8081]
    Basket --> Product
    Product --> DB[(PostgreSQL)]
    Basket --> Redis[(Redis)]
```
