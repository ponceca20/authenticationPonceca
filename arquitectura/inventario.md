::: mermaid
erDiagram
    %% --- Entidades Principales ---
    inv_ALMACEN {
        uint64 id PK
        uint64 empresa_id FK
        string nombre
        bool activo
    }

    inv_STOCK {
        uint64 empresa_id PK
        uint64 almacen_id PK
        uint64 producto_id PK
        decimal cantidad_actual
        timestamp updated_at
    }

    %% --- Movimientos (Simplificado) ---
    inv_MOVIMIENTO {
        uint64 id PK
        uint64 empresa_id FK
        uint64 usuario_id
       uint64 origen_operacion_id "ID de referencia a sistema externo, Nullable"
        string origen "ventas, compras, traslado, ajuste, etc."
        enum tipo "ENUM('ENTRADA','SALIDA','AJUSTE','TRASLADO')"        
        enum estado "ENUM('BORRADOR','PROCESADO','ANULADO')"       
        uint64 almacen_id FK "Para ENTRADA/SALIDA/AJUSTE"
        uint64 almacen_destino_id "Solo para TRASLADO"
        string descripcionmovimiento
        timestamp created_at
    }

    inv_MOVIMIENTO_DETALLE {
        uint64 id PK
        uint64 movimiento_id FK
        uint64 almacen_id FK "Index"
        uint64 almacen_destino_id FK "Almacén destino para este producto específico"
        uint64 producto_id
        decimal cantidad
        decimal saldo_anterior "Para auditoría"
        decimal saldo_resultante "Para auditoría"
        decimal costo_unitario
        uint64 lote_externo_id
        timestamp created_at
    }

    %% --- Relaciones ---
    inv_ALMACEN ||--|{ inv_STOCK : "tiene"
    inv_MOVIMIENTO ||--|{ inv_MOVIMIENTO_DETALLE : "contiene"
    inv_ALMACEN ||--o{ inv_MOVIMIENTO : "afecta_a"
:::
