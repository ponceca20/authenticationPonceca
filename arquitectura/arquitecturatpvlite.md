::: mermaid
graph TB
    %% Estilos
    classDef frontend fill:#e3f2fd,stroke:#2196f3,stroke-width:2px
    classDef gateway fill:#fce4ec,stroke:#e91e63,stroke-width:2px
    classDef service fill:#e8f5e9,stroke:#4caf50,stroke-width:2px
    classDef infrastructure fill:#fff3e0,stroke:#ff9800,stroke-width:2px
    classDef database fill:#f3e5f5,stroke:#9c27b0,stroke-width:2px
    classDef security fill:#efebe9,stroke:#795548,stroke-width:2px
    classDef observability fill:#e0f7fa,stroke:#00bcd4,stroke-width:2px

    subgraph Frontend ["Capa de Presentación"]
        UI[Frontend Angular]:::frontend
        MOBILE[App Móvil]:::frontend
    end

    subgraph Gateway ["API Gateway"]
        APIGW[API Gateway - Nginx]:::gateway
        SECURITY[Autenticación & Autorización]:::security
    end

    subgraph DomainServices ["Servicios de Dominio"]
        direction LR
        subgraph Ventas ["Dominio de Ventas"]
            SALES[Ventas]:::service
            PAYMENT[Pagos]:::service
            CUSTOMERS[Clientes]:::service
            DELIVERY[Delivery]:::service
        end

        subgraph Productos ["Dominio de Productos"]
            PRODUCTS[Productos]:::service
            OFFERS[Promociones]:::service
            RECIPES[Recetas]:::service
        end

        subgraph Inventario ["Dominio de Inventario"]
            INVENTORY[Inventario]:::service
            PURCHASE[Compras]:::service
            STOCKCONTROL[Control Stock]:::service
        end

        subgraph Operaciones ["Dominio de Operaciones"]
            PRODUCTION[Producción]:::service
            CLOSING[Caja]:::service
            RESERVATIONS[Reservas]:::service
        end
    end

    subgraph Core ["Servicios Core"]
        AUTH[IAM]:::service
        NOTIFICATION[Notificaciones]:::service
        REPORTS[Reportes]:::service
    end

    subgraph Infraestructura ["Infraestructura"]
        REDIS[Redis]:::infrastructure
        RABBITMQ[RabbitMQ]:::infrastructure
        MONGODB[MongoDB]:::infrastructure
        MYSQL[(MySQL)]:::database
    end

    subgraph Observabilidad ["Observabilidad"]
        PROMETHEUS[Prometheus]:::observability
        GRAFANA[Grafana]:::observability
        ELK[ELK Stack]:::observability
    end

    %% Conexiones principales
    Frontend --> APIGW
    APIGW --> SECURITY
    SECURITY --> DomainServices
    APIGW --> Core
    DomainServices <--> RABBITMQ
    DomainServices --> REDIS
    DomainServices --> MYSQL
    Core --> Infraestructura
    
    %% Observabilidad
    DomainServices --> Observabilidad
    Gateway --> Observabilidad
    Infraestructura --> Observabilidad
:::