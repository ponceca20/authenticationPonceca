::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ lot_LOTE : tiene
    PRODUCTO ||--o{ lot_LOTE : controla

    lot_LOTE {
        uint64 id PK "gorm:primaryKey"
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 producto_id FK "gorm:foreignKey"
        time fecha_vencimiento "null"
        decimal cantidad "precision(16,4)"
        bool activo "default:true"
        bool eliminado "default:false;index"
        bool requiere_vencimiento
        uint64 usuario_creacion_id FK
        uint64 usuario_actualizacion_id FK
    }
:::