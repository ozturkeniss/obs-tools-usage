# Microservices Architecture

## System Overview

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Microservices"
        Product[Product Service<br/>HTTP: 8080<br/>gRPC: 50050]
        Basket[Basket Service<br/>HTTP: 8081<br/>gRPC: 50051]
        Payment[Payment Service<br/>HTTP: 8082<br/>gRPC: 50052]
    end
    
    subgraph "Data Storage"
        PostgreSQL[(PostgreSQL<br/>Port: 5432)]
        Redis[(Redis<br/>Port: 6379)]
        MariaDB[(MariaDB<br/>Port: 3306)]
    end
    
    subgraph "Message Broker"
        Kafka[Apache Kafka<br/>Port: 9092<br/>JMX: 9101]
        Zookeeper[Zookeeper<br/>Port: 2181]
    end
    
    subgraph "Clients"
        HTTPClient[HTTP Client]
        GRPCClient[gRPC Client]
    end
    
    HTTPClient --> Product
    HTTPClient --> Basket
    HTTPClient --> Payment
    GRPCClient --> Product
    GRPCClient --> Basket
    GRPCClient --> Payment
    
    Product --> PostgreSQL
    Basket --> Redis
    Payment --> MariaDB
    
    Basket --> Product
    Payment --> Basket
    Payment --> Product
    
    Payment --> Kafka
    Kafka --> Zookeeper
```

## Product Service Architecture

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
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
    UseCase --> Repository
    UseCase --> DomainService
    
    Repository --> RepoImpl
    RepoImpl --> DB
    
    UseCase --> Config
    UseCase --> Logger
    UseCase --> Metrics
```

## Product Service API Endpoints

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph LR
    subgraph "HTTP Endpoints"
        GET1[GET /products<br/>Get all products]
        GET2[GET /products/{id}<br/>Get product by ID]
        POST[POST /products<br/>Create new product]
        PUT[PUT /products/{id}<br/>Update product]
        DELETE[DELETE /products/{id}<br/>Delete product]
        HEALTH[GET /health<br/>Health check]
    end
    
    subgraph "gRPC Methods"
        CreateProduct[CreateProduct]
        GetProduct[GetProduct]
        GetProducts[GetProducts]
        UpdateProduct[UpdateProduct]
        DeleteProduct[DeleteProduct]
    end
```

## Product Service Environment Variables

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Server Configuration"
        PORT[PORT: 8080]
        GRPC_PORT[GRPC_PORT: 50050]
        LOG_LEVEL[LOG_LEVEL: info]
    end
    
    subgraph "Database Configuration"
        DB_HOST[DB_HOST: localhost]
        DB_PORT[DB_PORT: 5432]
        DB_USER[DB_USER: postgres]
        DB_PASSWORD[DB_PASSWORD: password]
        DB_NAME[DB_NAME: product_service]
    end
    
    PORT --> GRPC_PORT
    GRPC_PORT --> LOG_LEVEL
    LOG_LEVEL --> DB_HOST
    DB_HOST --> DB_PORT
    DB_PORT --> DB_USER
    DB_USER --> DB_PASSWORD
    DB_PASSWORD --> DB_NAME
```

## Basket Service Architecture

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "External Layer"
        HTTP[HTTP API<br/>Port 8081]
        GRPC[gRPC API<br/>Port 50051]
        Redis[(Redis<br/>Port 6379)]
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
        Entity[Basket Entity]
        Repository[Repository Interface]
        ProductClient[Product Client]
    end
    
    subgraph "Infrastructure Layer"
        RedisImpl[Redis Implementation]
        GRPCClient[gRPC Product Client]
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
    UseCase --> Repository
    UseCase --> ProductClient
    
    Repository --> RedisImpl
    RedisImpl --> Redis
    
    ProductClient --> GRPCClient
    GRPCClient --> Product
    
    UseCase --> Config
    UseCase --> Logger
    UseCase --> Metrics
```

## Basket Service API Endpoints

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph LR
    subgraph "Basket Management"
        GET_BASKET[GET /baskets/{user_id}<br/>Get user basket]
        CREATE_BASKET[POST /baskets<br/>Create new basket]
        DELETE_BASKET[DELETE /baskets/{user_id}<br/>Delete basket]
    end
    
    subgraph "Item Management"
        ADD_ITEM[POST /baskets/{user_id}/items<br/>Add item]
        UPDATE_ITEM[PUT /baskets/{user_id}/items/{product_id}<br/>Update quantity]
        REMOVE_ITEM[DELETE /baskets/{user_id}/items/{product_id}<br/>Remove item]
        CLEAR_ITEMS[DELETE /baskets/{user_id}/items<br/>Clear all items]
    end
    
    subgraph "Health Check"
        HEALTH[GET /health<br/>Health check]
    end
    
    GET_BASKET --> ADD_ITEM
    ADD_ITEM --> UPDATE_ITEM
    UPDATE_ITEM --> REMOVE_ITEM
    REMOVE_ITEM --> CLEAR_ITEMS
    CLEAR_ITEMS --> HEALTH
```

## Basket Service Environment Variables

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Server Configuration"
        PORT[PORT: 8081]
        LOG_LEVEL[LOG_LEVEL: info]
    end
    
    subgraph "Redis Configuration"
        REDIS_HOST[REDIS_HOST: localhost]
        REDIS_PORT[REDIS_PORT: 6379]
        REDIS_PASSWORD[REDIS_PASSWORD: ]
        REDIS_DB[REDIS_DB: 0]
    end
    
    subgraph "Service Configuration"
        PRODUCT_SERVICE_URL[PRODUCT_SERVICE_URL: localhost:50050]
    end
    
    PORT --> LOG_LEVEL
    LOG_LEVEL --> REDIS_HOST
    REDIS_HOST --> REDIS_PORT
    REDIS_PORT --> REDIS_PASSWORD
    REDIS_PASSWORD --> REDIS_DB
    REDIS_DB --> PRODUCT_SERVICE_URL
```

## Payment Service Architecture

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "External Layer"
        HTTP[HTTP API<br/>Port 8082]
        GRPC[gRPC API<br/>Port 50052]
        MariaDB[(MariaDB<br/>Port 3306)]
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
        Entity[Payment Entity]
        Repository[Repository Interface]
        BasketClient[Basket Client]
        ProductClient[Product Client]
    end
    
    subgraph "Infrastructure Layer"
        MariaDBImpl[MariaDB Implementation]
        GRPCClient[gRPC Clients]
        Config[Configuration]
        Logger[Logging]
        Metrics[Prometheus Metrics]
        KafkaPublisher[Kafka Publisher]
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
    UseCase --> Repository
    UseCase --> BasketClient
    UseCase --> ProductClient
    
    Repository --> MariaDBImpl
    MariaDBImpl --> MariaDB
    
    BasketClient --> GRPCClient
    ProductClient --> GRPCClient
    GRPCClient --> Basket
    GRPCClient --> Product
    
    UseCase --> Config
    UseCase --> Logger
    UseCase --> Metrics
    UseCase --> KafkaPublisher
```

## Payment Service API Endpoints

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph LR
    subgraph "Payment Management"
        CREATE_PAYMENT[POST /payments<br/>Create payment]
        GET_PAYMENT[GET /payments/{id}<br/>Get payment]
        PROCESS_PAYMENT[POST /payments/{id}/process<br/>Process payment]
        CANCEL_PAYMENT[POST /payments/{id}/cancel<br/>Cancel payment]
        REFUND_PAYMENT[POST /payments/{id}/refund<br/>Refund payment]
    end
    
    subgraph "Payment History"
        GET_PAYMENTS[GET /payments<br/>Get all payments]
        GET_USER_PAYMENTS[GET /users/{user_id}/payments<br/>Get user payments]
    end
    
    subgraph "Health Check"
        HEALTH[GET /health<br/>Health check]
    end
    
    CREATE_PAYMENT --> GET_PAYMENT
    GET_PAYMENT --> PROCESS_PAYMENT
    PROCESS_PAYMENT --> CANCEL_PAYMENT
    CANCEL_PAYMENT --> REFUND_PAYMENT
    REFUND_PAYMENT --> GET_PAYMENTS
    GET_PAYMENTS --> GET_USER_PAYMENTS
    GET_USER_PAYMENTS --> HEALTH
```

## Payment Service Environment Variables

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Server Configuration"
        PORT[PORT: 8082]
        LOG_LEVEL[LOG_LEVEL: info]
    end
    
    subgraph "Database Configuration"
        DB_HOST[DB_HOST: localhost]
        DB_PORT[DB_PORT: 3306]
        DB_USER[DB_USER: payment]
        DB_PASSWORD[DB_PASSWORD: password]
        DB_NAME[DB_NAME: payment_service]
        DB_SSL_MODE[DB_SSL_MODE: false]
    end
    
    subgraph "Service Configuration"
        BASKET_SERVICE_URL[BASKET_SERVICE_URL: localhost:50051]
        PRODUCT_SERVICE_URL[PRODUCT_SERVICE_URL: localhost:50050]
    end
    
    PORT --> LOG_LEVEL
    LOG_LEVEL --> DB_HOST
    DB_HOST --> DB_PORT
    DB_PORT --> DB_USER
    DB_USER --> DB_PASSWORD
    DB_PASSWORD --> DB_NAME
    DB_NAME --> DB_SSL_MODE
    DB_SSL_MODE --> BASKET_SERVICE_URL
    BASKET_SERVICE_URL --> PRODUCT_SERVICE_URL
```

## Event-Driven Architecture with Kafka

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Event Publishers"
        PaymentService[Payment Service]
    end
    
    subgraph "Kafka Topics"
        PaymentEvents[payment-events]
        StockEvents[stock-events]
        BasketEvents[basket-events]
    end
    
    subgraph "Event Types"
        PaymentCompleted[Payment Completed]
        PaymentFailed[Payment Failed]
        PaymentRefunded[Payment Refunded]
        StockUpdated[Stock Updated]
        BasketCleared[Basket Cleared]
    end
    
    subgraph "Event Consumers"
        ProductService[Product Service]
        BasketService[Basket Service]
        NotificationService[Notification Service]
    end
    
    PaymentService --> PaymentEvents
    PaymentService --> StockEvents
    PaymentService --> BasketEvents
    
    PaymentEvents --> PaymentCompleted
    PaymentEvents --> PaymentFailed
    PaymentEvents --> PaymentRefunded
    
    StockEvents --> StockUpdated
    BasketEvents --> BasketCleared
    
    PaymentCompleted --> ProductService
    PaymentCompleted --> NotificationService
    PaymentFailed --> NotificationService
    PaymentRefunded --> NotificationService
    StockUpdated --> ProductService
    BasketCleared --> BasketService
```

## Docker Services Configuration

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Application Services"
        ProductService[product-service<br/>Ports: 8080, 50050]
        BasketService[basket-service<br/>Ports: 8081, 50051]
        PaymentService[payment-service<br/>Ports: 8082, 50052]
    end
    
    subgraph "Database Services"
        PostgreSQL[postgres<br/>Port: 5432]
        Redis[redis<br/>Port: 6379]
        MariaDB[mariadb<br/>Port: 3306]
    end
    
    subgraph "Message Broker Services"
        Kafka[kafka<br/>Ports: 9092, 9101]
        Zookeeper[zookeeper<br/>Port: 2181]
    end
    
    subgraph "Dependencies"
        ProductService --> PostgreSQL
        BasketService --> Redis
        PaymentService --> MariaDB
        BasketService --> ProductService
        PaymentService --> BasketService
        PaymentService --> ProductService
        PaymentService --> Kafka
        Kafka --> Zookeeper
    end
```

## Technology Stack

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Programming Language"
        Go[Go 1.22+]
    end
    
    subgraph "Frameworks"
        Gin[Gin HTTP Framework]
        GRPC[gRPC Framework]
        Wire[Wire Dependency Injection]
    end
    
    subgraph "Databases"
        PostgreSQL[PostgreSQL]
        Redis[Redis]
        MariaDB[MariaDB]
    end
    
    subgraph "Message Broker"
        Kafka[Apache Kafka]
        Zookeeper[Zookeeper]
    end
    
    subgraph "Monitoring"
        Prometheus[Prometheus Metrics]
        Logrus[Logrus Logging]
    end
    
    subgraph "Architecture Patterns"
        DDD[Domain Driven Design]
        CQRS[CQRS Pattern]
        CleanArch[Clean Architecture]
        EventDriven[Event-Driven Architecture]
    end
    
    Go --> Gin
    Go --> GRPC
    Go --> Wire
    
    Gin --> PostgreSQL
    GRPC --> Redis
    Wire --> MariaDB
    
    PostgreSQL --> Kafka
    Redis --> Zookeeper
    MariaDB --> Prometheus
    
    Kafka --> Logrus
    Zookeeper --> DDD
    
    Prometheus --> CQRS
    Logrus --> CleanArch
    DDD --> EventDriven
```

## Development Workflow

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph LR
    subgraph "Setup Phase"
        Setup[make setup]
        InstallDeps[make install-deps]
        Proto[make proto]
    end
    
    subgraph "Development Phase"
        Dev[make dev]
        Build[make build]
        Test[make test]
        Lint[make lint]
    end
    
    subgraph "Deployment Phase"
        DockerBuild[make docker-build]
        DockerRun[make docker-run]
        ServicesStart[make services-start]
    end
    
    subgraph "Maintenance Phase"
        ServicesStop[make services-stop]
        ServicesRestart[make services-restart]
        Clean[make clean]
    end
    
    Setup --> InstallDeps
    InstallDeps --> Proto
    Proto --> Dev
    Dev --> Build
    Build --> Test
    Test --> Lint
    Lint --> DockerBuild
    DockerBuild --> DockerRun
    DockerRun --> ServicesStart
    ServicesStart --> ServicesStop
    ServicesStop --> ServicesRestart
    ServicesRestart --> Clean
```

## Database Schema Overview

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
erDiagram
    PRODUCTS {
        int id PK
        string name
        string description
        decimal price
        int stock_quantity
        string category
        string image_url
        timestamp created_at
        timestamp updated_at
    }
    
    BASKETS {
        string id PK
        string user_id
        decimal total
        timestamp created_at
        timestamp updated_at
    }
    
    BASKET_ITEMS {
        string id PK
        string basket_id FK
        int product_id FK
        string name
        int quantity
        decimal price
        decimal subtotal
        string category
        timestamp created_at
    }
    
    PAYMENTS {
        string id PK
        string user_id
        string basket_id FK
        decimal amount
        string currency
        string status
        string method
        string provider
        string provider_id
        string description
        json metadata
        timestamp created_at
        timestamp updated_at
        timestamp processed_at
        timestamp expires_at
    }
    
    PAYMENT_ITEMS {
        string id PK
        string payment_id FK
        int product_id FK
        string name
        int quantity
        decimal price
        decimal subtotal
        string category
        timestamp created_at
    }
    
    BASKETS ||--o{ BASKET_ITEMS : contains
    PAYMENTS ||--o{ PAYMENT_ITEMS : contains
    PRODUCTS ||--o{ BASKET_ITEMS : referenced_by
    PRODUCTS ||--o{ PAYMENT_ITEMS : referenced_by
```

## API Request Flow

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
sequenceDiagram
    participant Client
    participant ProductService
    participant BasketService
    participant PaymentService
    participant Database
    participant Kafka
    
    Client->>ProductService: GET /products
    ProductService->>Database: Query products
    Database-->>ProductService: Return products
    ProductService-->>Client: Products response
    
    Client->>BasketService: POST /baskets/{user_id}/items
    BasketService->>ProductService: gRPC GetProduct
    ProductService-->>BasketService: Product details
    BasketService->>Database: Store basket item
    Database-->>BasketService: Item stored
    BasketService-->>Client: Item added
    
    Client->>PaymentService: POST /payments
    PaymentService->>BasketService: gRPC GetBasket
    BasketService-->>PaymentService: Basket details
    PaymentService->>Database: Create payment
    Database-->>PaymentService: Payment created
    PaymentService->>Kafka: Publish payment event
    PaymentService-->>Client: Payment response
```

## Service Communication Flow

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryText色': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Client Layer"
        WebClient[Web Client]
        MobileClient[Mobile Client]
        APIClient[API Client]
    end
    
    subgraph "API Gateway Layer"
        LoadBalancer[Load Balancer]
        RateLimiter[Rate Limiter]
        AuthGateway[Auth Gateway]
    end
    
    subgraph "Microservices Layer"
        ProductService[Product Service]
        BasketService[Basket Service]
        PaymentService[Payment Service]
    end
    
    subgraph "Data Layer"
        PostgreSQL[(PostgreSQL)]
        Redis[(Redis)]
        MariaDB[(MariaDB)]
    end
    
    subgraph "Message Layer"
        Kafka[Apache Kafka]
        EventStore[Event Store]
    end
    
    WebClient --> LoadBalancer
    MobileClient --> LoadBalancer
    APIClient --> LoadBalancer
    
    LoadBalancer --> RateLimiter
    RateLimiter --> AuthGateway
    
    AuthGateway --> ProductService
    AuthGateway --> BasketService
    AuthGateway --> PaymentService
    
    ProductService --> PostgreSQL
    BasketService --> Redis
    PaymentService --> MariaDB
    
    PaymentService --> Kafka
    Kafka --> EventStore
    
    BasketService --> ProductService
    PaymentService --> BasketService
    PaymentService --> ProductService
```