::: mermaid
erDiagram
    EMPRESA ||--o{ res_RESERVA : tiene
    res_RESERVA ||--o{ res_MESA_RESERVADA : contiene
    res_RESERVA }o--|| res_ESTADO_RESERVA : tiene
    res_RESERVA }o--|| CLIENTE : realiza
    res_RESERVA ||--o{ res_HISTORIAL_RESERVA : registra
    res_RESERVA ||--o{ res_DETALLE_RESERVA : incluye
    res_DETALLE_RESERVA }o--|| PRESENTACION : "reserva"

    res_RESERVA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 cliente_id FK
        uint64 sucursal_id FK "Sucursal donde se realiza la reserva"
        time fecha_reserva "Fecha y hora de la reserva"
        uint8 cantidad_personas
        string notas_especiales "Alergias, ocasión especial, etc."
        time fecha_creacion
        time fecha_actualizacion
        uint64 estado_id FK
        string codigo_confirmacion "Código único para el cliente"
        bool recordatorio_enviado
        uint64 usuario_creacion_id FK
        uint64 usuario_modificacion_id FK
        bool deleted
    }

    res_MESA_RESERVADA {
        uint64 id PK
        uint64 reserva_id FK
        uint64 mesa_id FK "Referencia a MESA del módulo Ventas"
        bool confirmada
    }

    res_ESTADO_RESERVA {
        uint64 id PK
        string codigo "PENDIENTE|CONFIRMADA|CANCELADA|COMPLETADA|NO_SHOW"
        string nombre
        string descripcion
        bool activo
    }

    res_HISTORIAL_RESERVA {
        uint64 id PK
        uint64 reserva_id FK
        uint64 estado_anterior_id FK
        uint64 estado_nuevo_id FK
        string comentario
        uint64 usuario_id FK
        time fecha_cambio
    }

    res_DETALLE_RESERVA {
        uint64 id PK
        uint64 reserva_id FK
        uint64 presentacion_id FK
        decimal cantidad "precision(16,4)"
        string notas_especiales
        decimal precio_unitario "precision(16,4)"
        decimal subtotal "precision(16,4)"
        bool confirmado
    }
::: 