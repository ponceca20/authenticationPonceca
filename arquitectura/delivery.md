::: mermaid
erDiagram
    EMPRESA ||--o{ del_PEDIDO_DELIVERY : tiene
    EMPRESA ||--o{ del_ZONA_ENTREGA : configura
    EMPRESA ||--o{ del_REPARTIDOR : emplea
    del_PEDIDO_DELIVERY ||--o{ del_TRACKING_PEDIDO : registra
    del_PEDIDO_DELIVERY }o--|| del_REPARTIDOR : asignado_a
    del_PEDIDO_DELIVERY }o--|| del_ZONA_ENTREGA : pertenece_a
    del_PEDIDO_DELIVERY }o--|| VENTA : corresponde_a
    del_REPARTIDOR ||--o{ del_DISPONIBILIDAD_REPARTIDOR : tiene

    del_PEDIDO_DELIVERY {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 venta_id FK "Referencia a la venta original"
        uint64 repartidor_id FK
        uint64 zona_entrega_id FK
        string direccion_entrega
        string referencias
        string telefono_contacto
        string nombre_contacto
        string estado "PENDIENTE|ASIGNADO|EN_CAMINO|ENTREGADO|CANCELADO"
        time fecha_pedido
        time fecha_asignacion
        time fecha_inicio_entrega
        time fecha_entrega_estimada
        time fecha_entrega_real
        float64 latitud_entrega
        float64 longitud_entrega
        decimal costo_delivery "precision(16,4)"
        string notas_repartidor
        string motivo_cancelacion
        bool deleted
    }

    del_TRACKING_PEDIDO {
        uint64 id PK
        uint64 pedido_delivery_id FK
        string estado
        time fecha_registro
        float64 latitud
        float64 longitud
        string comentario
        uint64 usuario_registro_id FK
    }

    del_REPARTIDOR {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 usuario_id FK "Referencia a usuario del sistema"
        string tipo_vehiculo "MOTO|BICICLETA|AUTO|PIE"
        string placa_vehiculo
        bool activo
        bool disponible
        time ultima_ubicacion_timestamp
        float64 ultima_latitud
        float64 ultima_longitud
        float64 calificacion_promedio
        int total_entregas
        bool deleted
    }

    del_ZONA_ENTREGA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        decimal costo_base "precision(16,4)"
        float64 radio_km
        bool activa
        json poligono "Coordenadas que definen la zona"
        time hora_inicio_atencion
        time hora_fin_atencion
        int tiempo_estimado_minutos
    }

    del_DISPONIBILIDAD_REPARTIDOR {
        uint64 id PK
        uint64 repartidor_id FK
        date fecha
        time hora_inicio
        time hora_fin
        string estado "DISPONIBLE|OCUPADO|DESCANSO"
        string notas
    }
::: 