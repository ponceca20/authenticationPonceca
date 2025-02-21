::: mermaid
erDiagram
    %% Relaciones principales
    
    CREDITO_SIMPLE ||--o{ PAGO : recibe
    
 

    CREDITO_SIMPLE {
        uint64 id PK
        uint64 venta_id FK
        uint64 cliente_id FK
        float64 monto
        time fecha_credito
        string estado "Pendiente, Pagado"
        string periodo "Diario, Semanal, Quincenal"
        time fecha_limite_pago
    }

    PAGO {
        uint64 id PK
        uint64 credito_id FK
        float64 monto
        time fecha_pago
        string forma_pago
        string referencia
    }

    CLIENTE_CREDITO {
        uint64 id PK
        uint64 cliente_id FK
        float64 limite_credito
        string periodo_pago
        bool activo
        float64 saldo_actual
        time ultima_fecha_pago
    }
:::
