::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ pay_DOCUMENTO : tiene
    EMPRESA ||--o{ pay_PAGO : tiene
    EMPRESA ||--o{ pay_RECORDATORIO_PAGO : tiene
    pay_DOCUMENTO ||--o{ pay_PAGO : tiene
    pay_PAGO }o--|| pay_METODO_PAGO : utiliza
    pay_DOCUMENTO }o--|| pay_ESTADO_CREDITO : tiene
    pay_PAGO }o--|| pay_ESTADO_PAGO : tiene
    pay_DOCUMENTO ||--o{ pay_RECORDATORIO_PAGO : tiene

    pay_DOCUMENTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string tipo "COMPRA|VENTA"
        decimal monto_total "precision(16,4)"
        decimal monto_pendiente "precision(16,4)"
        time fecha_emision
        time fecha_vencimiento
        bool es_credito
        uint64 documento_origen_id FK
        string tipo_origen "VENTA|COMPRA"
        uint64 cliente_proveedor_id FK
        string serie_origen
        string numero_origen
        uint64 estado_pago_id FK
    }

    pay_PAGO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 documento_id FK
        uint64 metodo_pago_id FK
        decimal monto "precision(16,4)"
        string referencia
        time fecha_pago
        bool confirmado
        string estado_validacion
        time fecha_validacion
        string codigo_operacion
        uint64 usuario_validacion_id FK
    }

    pay_ESTADO_CREDITO {
        uint64 id PK
        string codigo
        string nombre
        bool activo
    }

    pay_ESTADO_PAGO {
        uint64 id PK
        string codigo "PENDIENTE|PARCIAL|PAGADO|VENCIDO"
        string nombre
        bool activo
    }

    pay_METODO_PAGO {
        uint64 id PK
        string codigo
        string nombre
        bool activo
        bool requiere_validacion "Para Yape, Plin, transferencias, etc."
    }

    pay_RECORDATORIO_PAGO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 documento_id FK
        time fecha_vencimiento
        bool notificado
        string correo_proveedor
        string telefono_proveedor
        uint64 usuario_creacion_id FK
        uint64 usuario_actualizacion_id FK
    }
:::
