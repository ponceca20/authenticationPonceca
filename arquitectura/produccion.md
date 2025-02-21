::: mermaid
erDiagram
    %% Nuevas relaciones para soporte multi-empresa
    EMPRESA ||--o{ prd_AREA : tiene
    EMPRESA ||--o{ prd_RECETA : tiene
    EMPRESA ||--o{ prd_PEDIDO_PRODUCCION : tiene
    EMPRESA ||--o{ prd_MERMA : tiene
    EMPRESA ||--o{ prd_TRANSFORMACION : tiene

    %% Resto de relaciones existentes
    prd_AREA ||--o{ prd_PEDIDO_PRODUCCION : solicita
    prd_PEDIDO_PRODUCCION }o--|| ALMACEN : despacha_desde
    prd_RECETA ||--o{ prd_DETALLE_RECETA : contiene
    prd_RECETA }o--|| PRODUCTO : pertenece
    prd_DETALLE_RECETA }o--|| PRESENTACION : usa
    prd_PEDIDO_PRODUCCION ||--o{ prd_DETALLE_PEDIDO_PRODUCCION : contiene
    prd_DETALLE_PEDIDO_PRODUCCION }o--|| PRESENTACION : requiere
    prd_MERMA ||--o{ prd_DETALLE_MERMA : contiene
    prd_DETALLE_MERMA }o--|| PRESENTACION : afecta_a
    prd_TRANSFORMACION ||--o{ prd_DETALLE_TRANSFORMACION : contiene
    prd_DETALLE_TRANSFORMACION }o--|| PRESENTACION : utiliza

    %% Entidades actualizadas
    prd_AREA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre "Parrilla, Pastas, Ensaladas, etc"
        string tipo "Producción, Almacén"
        bool activo
        uint64 responsable_id FK
    }

    prd_RECETA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 producto_id FK
        string nombre
        string descripcion
        decimal porciones_receta "precision(16,4)"
        string instrucciones_preparacion
        string notas_especiales
        decimal costo_estimado "precision(16,4)"
        decimal tiempo_preparacion "precision(16,4)"
        bool activo
        uint64 usuario_creacion_id FK
        uint64 usuario_modificacion_id FK
    }

    prd_DETALLE_RECETA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 receta_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        string nombre_insumo
        decimal cantidad "precision(16,4)"
        decimal merma_estimada "precision(16,4)"
        string unidad_medida
        string notas
        bool es_opcional
        bool puede_sustituirse
        string sustitutos_permitidos
    }

    prd_PEDIDO_PRODUCCION {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 area_solicitante_id FK
        uint64 almacen_despacho_id FK "Almacén que despachará los insumos"
        string numero_pedido
        time fecha_pedido
        time fecha_requerida
        string estado "Pendiente, En Proceso, Completado, Cancelado"
        string prioridad "Normal, Urgente"
        string notas
        uint usuario_solicitud_id FK
        uint64 usuario_despacho_id FK
        time fecha_despacho
        bool deleted
    }

    prd_DETALLE_PEDIDO_PRODUCCION {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 pedido_produccion_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        string nombre_producto
        decimal cantidad_solicitada "precision(16,4)"
        decimal cantidad_despachada "precision(16,4)"
        string estado "Pendiente, En Proceso, Completado"
        string notas
        time hora_inicio_preparacion
        time hora_fin_preparacion
        uint64 usuario_preparacion_id FK
        bool deleted
    }

    prd_MERMA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string numero_documento
        time fecha_registro
        string tipo_merma "Producción, Manipulación, Vencimiento, Caducidad"
        string estado "Pendiente, Aprobado, Rechazado"
        string motivo_general
        uint64 usuario_registro_id FK
        uint64 usuario_aprobacion_id FK "Null hasta aprobación"
        time fecha_aprobacion
        string notas
        bool deleted
    }

    prd_DETALLE_MERMA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 merma_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        string nombre_producto
        decimal cantidad "precision(16,4)"
        string unidad_medida
        decimal costo_unitario "precision(16,4)"
        decimal costo_total "precision(16,4)"
        string motivo_especifico
        string lote "Opcional, para control de lotes"
        time fecha_vencimiento "Opcional"
        string notas
        bool deleted
    }

    prd_TRANSFORMACION {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string numero_documento
        time fecha_registro
        string estado "Pendiente, Completado, Anulado"
        uint64 presentacion_origen_id FK "Insumo a transformar"
        uint64 producto_origen_id FK "Producto a transformar"
        decimal cantidad_origen "precision(16,4)"
        string unidad_medida_origen
        uint64 usuario_registro_id FK
        bool procesado "Indica si ya se procesó el movimiento de stock"
        time fecha_procesado "Fecha cuando se procesó el movimiento"
        uint64 usuario_procesado_id FK "Usuario que procesó el movimiento"
        bool revertido "Indica si se revirtió el procesamiento"
        time fecha_reversion "Fecha cuando se revirtió"
        uint64 usuario_reversion_id FK "Usuario que revirtió"
        string motivo_reversion
        bool deleted
    }

    prd_DETALLE_TRANSFORMACION {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 transformacion_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        decimal cantidad "precision(16,4)"
        string unidad_medida
        string lote
        decimal costo_unitario "precision(16,4)"
        bool deleted
    }
::: 