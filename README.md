# Microservices Architecture

## System Overview

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "API Gateway"
        Gateway[FiberV2 Gateway<br/>HTTP: 8083<br/>Load Balancer<br/>Circuit Breaker<br/>Rate Limiting]
    end
    
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
    
    HTTPClient --> Gateway
    GRPCClient --> Gateway
    
    Gateway --> Product
    Gateway --> Basket
    Gateway --> Payment
    
    Product --> PostgreSQL
    Basket --> Redis
    Payment --> MariaDB
    
    Basket --> Product
    Payment --> Basket
    Payment --> Product
    
    Payment --> Kafka
    Kafka --> Zookeeper
```

## FiberV2 Gateway Architecture

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "External Layer"
        HTTPClient[HTTP Client]
        AdminClient[Admin Client]
    end
    
    subgraph "Gateway Layer"
        Gateway[FiberV2 Gateway<br/>Port: 8083]
        Middleware[Middleware<br/>CORS, Logging, Metrics<br/>Rate Limiting, Security]
    end
    
    subgraph "Routing Layer"
        Router[Router<br/>Service Routing<br/>Path Rewriting]
        LoadBalancer[Load Balancer<br/>Round Robin<br/>Least Connections<br/>Weighted Round Robin]
    end
    
    subgraph "Circuit Breaker Layer"
        CircuitBreaker[Circuit Breaker<br/>Failure Detection<br/>Service Isolation<br/>Auto Recovery]
    end
    
    subgraph "Proxy Layer"
        ReverseProxy[Reverse Proxy<br/>Request Forwarding<br/>Response Handling<br/>Header Management]
    end
    
    subgraph "Backend Services"
        ProductService[Product Service<br/>Port: 8080]
        BasketService[Basket Service<br/>Port: 8081]
        PaymentService[Payment Service<br/>Port: 8082]
    end
    
    HTTPClient --> Gateway
    AdminClient --> Gateway
    
    Gateway --> Middleware
    Middleware --> Router
    
    Router --> LoadBalancer
    LoadBalancer --> CircuitBreaker
    
    CircuitBreaker --> ReverseProxy
    ReverseProxy --> ProductService
    ReverseProxy --> BasketService
    ReverseProxy --> PaymentService
```

## FiberV2 Gateway Features

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Core Features"
        ReverseProxy[Reverse Proxy<br/>Request Forwarding<br/>Response Handling]
        LoadBalancing[Load Balancing<br/>Round Robin<br/>Least Connections<br/>Weighted Round Robin]
        CircuitBreaker[Circuit Breaker<br/>Failure Detection<br/>Service Isolation<br/>Auto Recovery]
    end
    
    subgraph "Security Features"
        RateLimiting[Rate Limiting<br/>Request Throttling<br/>Burst Control]
        CORSSupport[CORS Support<br/>Cross-Origin Requests<br/>Header Management]
        SecurityHeaders[Security Headers<br/>XSS Protection<br/>CSRF Protection<br/>HSTS]
    end
    
    subgraph "Monitoring Features"
        HealthChecks[Health Checks<br/>Service Monitoring<br/>Status Reporting]
        Metrics[Prometheus Metrics<br/>Request Counters<br/>Response Times<br/>Error Rates]
        Logging[Structured Logging<br/>Request Tracking<br/>Error Logging]
    end
    
    subgraph "Admin Features"
        AdminAPI[Admin API<br/>Service Management<br/>Configuration Updates]
        StatusMonitoring[Status Monitoring<br/>Real-time Health<br/>Performance Metrics]
        ServiceDiscovery[Service Discovery<br/>Dynamic Backend<br/>Configuration]
    end
    
    ReverseProxy --> LoadBalancing
    LoadBalancing --> CircuitBreaker
    CircuitBreaker --> RateLimiting
    RateLimiting --> CORSSupport
    CORSSupport --> SecurityHeaders
    SecurityHeaders --> HealthChecks
    HealthChecks --> Metrics
    Metrics --> Logging
    Logging --> AdminAPI
    AdminAPI --> StatusMonitoring
    StatusMonitoring --> ServiceDiscovery
```

## FiberV2 Gateway API Endpoints

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph LR
    subgraph "Service Endpoints"
        ProductAPI[GET /api/products/*<br/>Product Service Proxy]
        BasketAPI[GET /api/baskets/*<br/>Basket Service Proxy]
        PaymentAPI[GET /api/payments/*<br/>Payment Service Proxy]
    end
    
    subgraph "Admin Endpoints"
        GatewayStatus[GET /admin/status<br/>Gateway Status]
        ServiceStatus[GET /admin/services<br/>Service Status]
        LoadBalancerStats[GET /admin/loadbalancer/:service<br/>Load Balancer Stats]
        CircuitBreakerStats[GET /admin/circuitbreaker/:service<br/>Circuit Breaker Stats]
    end
    
    subgraph "Health Endpoints"
        HealthCheck[GET /health<br/>Health Check]
        DetailedHealth[GET /health/detailed<br/>Detailed Health Check]
        ReadinessCheck[GET /health/ready<br/>Readiness Check]
        LivenessCheck[GET /health/live<br/>Liveness Check]
    end
    
    subgraph "Metrics Endpoint"
        Metrics[GET /metrics<br/>Prometheus Metrics]
    end
    
    ProductAPI --> GatewayStatus
    BasketAPI --> ServiceStatus
    PaymentAPI --> LoadBalancerStats
    GatewayStatus --> CircuitBreakerStats
    ServiceStatus --> HealthCheck
    LoadBalancerStats --> DetailedHealth
    CircuitBreakerStats --> ReadinessCheck
    HealthCheck --> LivenessCheck
    DetailedHealth --> Metrics
```

## FiberV2 Gateway Environment Variables

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Server Configuration"
        PORT[PORT: 8080]
        LOG_LEVEL[LOG_LEVEL: info]
        LOG_FORMAT[LOG_FORMAT: json]
    end
    
    subgraph "Service Configuration"
        PRODUCT_ENABLED[PRODUCT_SERVICE_ENABLED: true]
        PRODUCT_URLS[PRODUCT_SERVICE_URLS: http://product-service:8080]
        BASKET_ENABLED[BASKET_SERVICE_ENABLED: true]
        BASKET_URLS[BASKET_SERVICE_URLS: http://basket-service:8081]
        PAYMENT_ENABLED[PAYMENT_SERVICE_ENABLED: true]
        PAYMENT_URLS[PAYMENT_SERVICE_URLS: http://payment-service:8082]
    end
    
    subgraph "Circuit Breaker Configuration"
        CB_ENABLED[CIRCUIT_BREAKER_ENABLED: true]
        CB_MAX_REQUESTS[CIRCUIT_BREAKER_MAX_REQUESTS: 10]
        CB_INTERVAL[CIRCUIT_BREAKER_INTERVAL: 60]
        CB_TIMEOUT[CIRCUIT_BREAKER_TIMEOUT: 30]
    end
    
    subgraph "Load Balancer Configuration"
        LB_ENABLED[LOAD_BALANCER_ENABLED: true]
        LB_STRATEGY[LOAD_BALANCER_STRATEGY: round_robin]
    end
    
    subgraph "Rate Limiting Configuration"
        RL_ENABLED[RATE_LIMIT_ENABLED: true]
        RL_REQUESTS[RATE_LIMIT_REQUESTS: 100]
        RL_WINDOW[RATE_LIMIT_WINDOW: 1m]
        RL_BURST[RATE_LIMIT_BURST: 10]
    end
    
    PORT --> LOG_LEVEL
    LOG_LEVEL --> LOG_FORMAT
    LOG_FORMAT --> PRODUCT_ENABLED
    PRODUCT_ENABLED --> PRODUCT_URLS
    PRODUCT_URLS --> BASKET_ENABLED
    BASKET_ENABLED --> BASKET_URLS
    BASKET_URLS --> PAYMENT_ENABLED
    PAYMENT_ENABLED --> PAYMENT_URLS
    PAYMENT_URLS --> CB_ENABLED
    CB_ENABLED --> CB_MAX_REQUESTS
    CB_MAX_REQUESTS --> CB_INTERVAL
    CB_INTERVAL --> CB_TIMEOUT
    CB_TIMEOUT --> LB_ENABLED
    LB_ENABLED --> LB_STRATEGY
    LB_STRATEGY --> RL_ENABLED
    RL_ENABLED --> RL_REQUESTS
    RL_REQUESTS --> RL_WINDOW
    RL_WINDOW --> RL_BURST
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
    subgraph "Gateway Services"
        Gateway[fiberv2-gateway<br/>Port: 8083<br/>Load Balancer<br/>Circuit Breaker<br/>Rate Limiting]
    end
    
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
        Gateway --> ProductService
        Gateway --> BasketService
        Gateway --> PaymentService
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
    
    subgraph "Gateway Frameworks"
        Fiber[Fiber HTTP Framework]
        FastHTTP[FastHTTP]
        CircuitBreaker[Circuit Breaker]
        LoadBalancer[Load Balancer]
    end
    
    subgraph "Microservice Frameworks"
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
        GatewayPattern[API Gateway Pattern]
    end
    
    Go --> Fiber
    Go --> Gin
    Go --> GRPC
    Go --> Wire
    
    Fiber --> FastHTTP
    Fiber --> CircuitBreaker
    Fiber --> LoadBalancer
    
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
    EventDriven --> GatewayPattern
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
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Client Layer"
        WebClient[Web Client]
        MobileClient[Mobile Client]
        APIClient[API Client]
    end
    
    subgraph "API Gateway Layer"
        FiberGateway[FiberV2 Gateway<br/>Port: 8083]
        LoadBalancer[Load Balancer<br/>Round Robin<br/>Least Connections]
        CircuitBreaker[Circuit Breaker<br/>Failure Detection<br/>Service Isolation]
        RateLimiter[Rate Limiter<br/>Request Throttling<br/>Burst Control]
        ReverseProxy[Reverse Proxy<br/>Request Forwarding<br/>Response Handling]
    end
    
    subgraph "Microservices Layer"
        ProductService[Product Service<br/>Port: 8080]
        BasketService[Basket Service<br/>Port: 8081]
        PaymentService[Payment Service<br/>Port: 8082]
    end
    
    subgraph "Data Layer"
        PostgreSQL[(PostgreSQL<br/>Port: 5432)]
        Redis[(Redis<br/>Port: 6379)]
        MariaDB[(MariaDB<br/>Port: 3306)]
    end
    
    subgraph "Message Layer"
        Kafka[Apache Kafka<br/>Port: 9092]
        Zookeeper[Zookeeper<br/>Port: 2181]
    end
    
    WebClient --> FiberGateway
    MobileClient --> FiberGateway
    APIClient --> FiberGateway
    
    FiberGateway --> LoadBalancer
    LoadBalancer --> CircuitBreaker
    CircuitBreaker --> RateLimiter
    RateLimiter --> ReverseProxy
    
    ReverseProxy --> ProductService
    ReverseProxy --> BasketService
    ReverseProxy --> PaymentService
    
    ProductService --> PostgreSQL
    BasketService --> Redis
    PaymentService --> MariaDB
    
    PaymentService --> Kafka
    Kafka --> Zookeeper
    
    BasketService --> ProductService
    PaymentService --> BasketService
    PaymentService --> ProductService
```

## Kubernetes Deployment Architecture

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Kubernetes Cluster"
        subgraph "Ingress Layer"
            Ingress[Ingress Controller<br/>nginx-ingress<br/>SSL Termination<br/>Load Balancing]
        end
        
        subgraph "Application Layer"
            Gateway[Gateway Deployment<br/>3 Replicas<br/>LoadBalancer Service]
            ProductService[Product Service<br/>2 Replicas<br/>ClusterIP Service]
            BasketService[Basket Service<br/>2 Replicas<br/>ClusterIP Service]
            PaymentService[Payment Service<br/>2 Replicas<br/>ClusterIP Service]
        end
        
        subgraph "Data Layer"
            PostgreSQL[PostgreSQL<br/>Bitnami Chart<br/>8Gi Persistent Volume]
            Redis[Redis<br/>Bitnami Chart<br/>4Gi Persistent Volume]
            MariaDB[MariaDB<br/>Bitnami Chart<br/>8Gi Persistent Volume]
        end
        
        subgraph "Message Layer"
            Kafka[Kafka<br/>Bitnami Chart<br/>3 Replicas<br/>10Gi Persistent Volume]
            Zookeeper[Zookeeper<br/>Bitnami Chart<br/>3 Replicas]
        end
        
        subgraph "Monitoring Layer"
            ServiceMonitor[ServiceMonitor<br/>Prometheus Integration]
            HPA[Horizontal Pod Autoscaler<br/>CPU/Memory Based Scaling]
        end
        
        subgraph "Security Layer"
            NetworkPolicy[NetworkPolicy<br/>Traffic Isolation]
            ServiceAccount[ServiceAccount<br/>RBAC Integration]
            PodSecurityContext[Pod Security Context<br/>Non-root User<br/>Read-only Filesystem]
        end
    end
    
    subgraph "External Access"
        Client[External Client]
        LoadBalancer[Load Balancer<br/>AWS NLB / GCP LB]
    end
    
    Client --> LoadBalancer
    LoadBalancer --> Ingress
    Ingress --> Gateway
    
    Gateway --> ProductService
    Gateway --> BasketService
    Gateway --> PaymentService
    
    ProductService --> PostgreSQL
    BasketService --> Redis
    PaymentService --> MariaDB
    
    PaymentService --> Kafka
    Kafka --> Zookeeper
    
    BasketService --> ProductService
    PaymentService --> BasketService
    PaymentService --> ProductService
    
    Gateway --> ServiceMonitor
    Gateway --> HPA
    
    Gateway --> NetworkPolicy
    ProductService --> NetworkPolicy
    BasketService --> NetworkPolicy
    PaymentService --> NetworkPolicy
    
    Gateway --> ServiceAccount
    ProductService --> ServiceAccount
    BasketService --> ServiceAccount
    PaymentService --> ServiceAccount
```

## Helm Chart Structure

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Helm Chart"
        ChartYaml[Chart.yaml<br/>Metadata & Dependencies]
        ValuesYaml[values.yaml<br/>Configuration Values]
        
        subgraph "Templates"
            Deployments[Deployments<br/>product-service<br/>basket-service<br/>payment-service<br/>gateway]
            Services[Services<br/>ClusterIP Services<br/>LoadBalancer Service]
            ConfigMaps[ConfigMaps<br/>Service Configurations]
            Ingress[Ingress<br/>External Access]
            ServiceAccount[ServiceAccount<br/>RBAC]
            NetworkPolicy[NetworkPolicy<br/>Security]
            ServiceMonitor[ServiceMonitor<br/>Monitoring]
            HPA[HPA<br/>Autoscaling]
        end
        
        subgraph "Dependencies"
            PostgreSQLChart[PostgreSQL Chart<br/>Bitnami]
            RedisChart[Redis Chart<br/>Bitnami]
            MariaDBChart[MariaDB Chart<br/>Bitnami]
            KafkaChart[Kafka Chart<br/>Bitnami]
        end
    end
    
    ChartYaml --> Deployments
    ValuesYaml --> Deployments
    
    Deployments --> Services
    Services --> ConfigMaps
    ConfigMaps --> Ingress
    Ingress --> ServiceAccount
    ServiceAccount --> NetworkPolicy
    NetworkPolicy --> ServiceMonitor
    ServiceMonitor --> HPA
    
    ChartYaml --> PostgreSQLChart
    ChartYaml --> RedisChart
    ChartYaml --> MariaDBChart
    ChartYaml --> KafkaChart
```

## Deployment Commands

```bash
# Install the Helm chart
helm install obs-tools-usage ./helm

# Upgrade the deployment
helm upgrade obs-tools-usage ./helm

# Check deployment status
helm status obs-tools-usage

# View all resources
kubectl get all -l app.kubernetes.io/name=obs-tools-usage

# Access the gateway
kubectl port-forward svc/obs-tools-usage-gateway 8080:8080

# View logs
kubectl logs -l app.kubernetes.io/component=gateway -f

# Scale services
kubectl scale deployment obs-tools-usage-gateway --replicas=5
```

## Environment Configuration

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Production Environment"
        ProdValues[values.yaml<br/>Production Settings]
        ProdSecurity[Security Context<br/>Non-root User<br/>Read-only Filesystem]
        ProdResources[Resource Limits<br/>CPU: 1000m<br/>Memory: 1Gi]
        ProdPersistence[Persistent Volumes<br/>8Gi PostgreSQL<br/>4Gi Redis<br/>8Gi MariaDB<br/>10Gi Kafka]
    end
    
    subgraph "Development Environment"
        DevValues[values-dev.yaml<br/>Development Settings]
        DevSecurity[Relaxed Security<br/>Debug Mode]
        DevResources[Lower Resources<br/>CPU: 250m<br/>Memory: 256Mi]
        DevPersistence[Smaller Volumes<br/>1Gi each]
    end
    
    subgraph "Staging Environment"
        StagingValues[values-staging.yaml<br/>Staging Settings]
        StagingSecurity[Production-like Security]
        StagingResources[Medium Resources<br/>CPU: 500m<br/>Memory: 512Mi]
        StagingPersistence[Medium Volumes<br/>4Gi each]
    end
```

## Monitoring and Observability

```mermaid
%%{init: {'theme':'base', 'themeVariables': { 'primaryColor': '#663399', 'primaryTextColor': '#ffffff', 'primaryBorderColor': '#663399', 'lineColor': '#ffffff', 'secondaryColor': '#663399', 'tertiaryColor': '#663399'}}}%%
graph TB
    subgraph "Application Metrics"
        PrometheusMetrics[Prometheus Metrics<br/>Request Count<br/>Response Time<br/>Error Rate<br/>CPU Usage<br/>Memory Usage]
    end
    
    subgraph "Health Checks"
        LivenessProbe[Liveness Probe<br/>/health endpoint<br/>30s interval]
        ReadinessProbe[Readiness Probe<br/>/health/ready endpoint<br/>5s interval]
    end
    
    subgraph "Logging"
        StructuredLogs[Structured Logging<br/>JSON Format<br/>Request Tracking<br/>Error Logging]
    end
    
    subgraph "Service Discovery"
        ServiceMonitor[ServiceMonitor<br/>Prometheus Integration<br/>30s scrape interval]
    end
    
    subgraph "Autoscaling"
        HPA[HPA<br/>CPU-based Scaling<br/>Memory-based Scaling<br/>Min: 1, Max: 100]
    end
    
    PrometheusMetrics --> LivenessProbe
    LivenessProbe --> ReadinessProbe
    ReadinessProbe --> StructuredLogs
    StructuredLogs --> ServiceMonitor
    ServiceMonitor --> HPA
```