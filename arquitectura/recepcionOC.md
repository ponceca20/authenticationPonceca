::: mermaid
erDiagram
    EMPRESA ||--o{ rec_RECEPCION_OC : tiene
    rec_RECEPCION_OC ||--o{ rec_DETALLE_RECEPCION_OC : contiene
    rec_EQUIVALENCIA_PRODUCTO ||--o{ rec_DETALLE_RECEPCION_OC : utiliza
    EMPRESA ||--o{ rec_EQUIVALENCIA_PRODUCTO : tiene

    rec_RECEPCION_OC {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 cliente_id FK "Referencia a Cliente en MS Ventas"
        string numero_oc_cliente UK
        time fecha_recepcion
        string estado "PENDIENTE/APROBADO/RECHAZADO"
        boolean activo
        bool deleted
    }

    rec_DETALLE_RECEPCION_OC {
        uint64 id PK
        uint64 recepcion_oc_id FK
        uint64 producto_id FK "Referencia a Producto en MS Productos"
        uint64 equivalencia_id FK
        string codigo_producto_cliente "Código original del cliente"
        decimal cantidad_solicitada "precision(16,4)"
        decimal cantidad_aceptada "precision(16,4)"
        decimal precio_acordado "precision(16,4)"
        decimal subtotal_aceptado "precision(16,4)"
        string estado "PENDIENTE/APROBADO/RECHAZADO"
        string motivo_rechazo
    }

    rec_EQUIVALENCIA_PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 cliente_id FK "Referencia a Cliente en MS Ventas"
        uint64 producto_id FK "Referencia a Producto en MS Productos"
        string codigo_producto_cliente UK
        string descripcion_cliente
        boolean activo
        time ultima_actualizacion
        bool deleted
    }
:::