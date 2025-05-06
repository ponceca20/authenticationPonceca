::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ lot_LOTE : tiene
    PRODUCTO ||--o{ lot_LOTE : controla
    lot_LOTE ||--o{ lot_MOVIMIENTO : registra
    lot_LOTE ||--o{ lot_ALERTA : configura
    lot_LOTE ||--o{ lot_DOCUMENTO : tiene

    lot_LOTE {
        uint64 id PK
        uint64 empresa_id FK
        uint64 producto_id FK
        string codigo_lote UK
        string codigo_lote_proveedor
        uint64 proveedor_id FK
        time fecha_fabricacion
        time fecha_vencimiento
        decimal cantidad_inicial "precision(16,4)"
        decimal cantidad_actual "precision(16,4)"
        decimal costo_unitario "precision(16,4)"
        string estado "VIGENTE|POR_VENCER|VENCIDO|AGOTADO"
        bool activo
        bool eliminado
        bool requiere_vencimiento
        uint64 usuario_creacion_id FK
        uint64 usuario_actualizacion_id FK
        timestamp created_at
        timestamp updated_at
    }

    lot_MOVIMIENTO {
        uint64 id PK
        uint64 lote_id FK
        string tipo "ENTRADA|SALIDA|AJUSTE|MERMA"
        decimal cantidad "precision(16,4)"
        decimal costo_unitario "precision(16,4)"
        string motivo 
        string documento_referencia
        uint64 usuario_id FK
        timestamp fecha_movimiento
        string notas
    }

    lot_ALERTA {
        uint64 id PK
        uint64 lote_id FK
        int dias_anticipacion
        string tipo_alerta "VENCIMIENTO|STOCK_MINIMO"
        string estado "PENDIENTE|NOTIFICADO|ATENDIDO"
        bool notificado
        timestamp fecha_alerta
        timestamp fecha_notificacion
        string destinatarios "Correos/usuarios a notificar"
    }

    lot_DOCUMENTO {
        uint64 id PK
        uint64 lote_id FK
        string tipo "FACTURA|GUIA|CERTIFICADO"
        string numero_documento
        string url_documento
        timestamp fecha_documento
        string notas
    }
:::