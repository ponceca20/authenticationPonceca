::: mermaid
erDiagram
    EMPRESA ||--o{ stk_PREPARACION_DIARIA : tiene
    
    stk_PREPARACION_DIARIA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 producto_id FK "Referencia a PRODUCTO del microservicio productos"
        datetime inicio_jornada "Timestamp inicio del día operativo"
        datetime fin_jornada "Timestamp fin del día operativo"
        decimal cantidad_preparada "precision(16,4)"
        decimal cantidad_vendida "precision(16,4)"
        bool jornada_activa "Indica si la jornada está en curso"
    }
::: 