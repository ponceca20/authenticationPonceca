::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ com_COMPRA : tiene
    EMPRESA ||--o{ com_PROVEEDOR : tiene
    EMPRESA ||--o{ com_TIPO_DOCUMENTO_COMPRAS : tiene
    EMPRESA ||--o{ com_ORDEN_COMPRA : tiene
    com_COMPRA ||--o{ com_DETALLE_COMPRA : contiene
    com_COMPRA }o--|| com_PROVEEDOR : realiza
    com_COMPRA }|--|| com_ESTADO_DOCUMENTO : tiene
    com_COMPRA ||--o{ com_IMAGEN_DOCUMENTO_COMPRA : tiene
    com_ORDEN_COMPRA ||--o{ com_DETALLE_ORDEN_COMPRA : contiene
    com_ORDEN_COMPRA }o--|| com_PROVEEDOR : realiza
    com_ORDEN_COMPRA }|--|| com_ESTADO_ORDEN : tiene
    com_DETALLE_COMPRA ||--o{ com_DETALLE_RECEPCION : tiene
    com_RECEPCION ||--o{ com_DETALLE_RECEPCION : contiene
    com_DETALLE_RECEPCION ||--o{ com_DETALLE_RECEPCION_LOTE : contiene
    com_PROVEEDOR ||--o{ com_PRODUCTO_PROVEEDOR : tiene
    com_PRODUCTO_PROVEEDOR ||--o{ com_HISTORIAL_PRECIOS_PROVEEDOR : registra

    %% Entidad COMPRA
    com_COMPRA {
        uint64 id PK
        uint64 empresa_id FK
        uint64 proveedor_id FK
        uint64 tipo_documento_compras_id FK
        uint64 estado_documento_id FK
        string serie
        string numero
        decimal subtotal "precision(16,4)"
        decimal impuestos "precision(16,4)"
        decimal total "precision(16,4)"
        bool activo
        time fecha_emision
        time fecha_recepcion
        decimal descuento_porcentual "precision(16,4)"
        decimal descuento_importe "precision(16,4)"
        string guia_remision
        bool facturado
        bool eliminado
        string orden_compra_ref
        uint64 usuario_creacion_id FK
        bool requiere_credito
        uint64 documento_credito_id FK
        string tipo_comprobante_electronico
        string identificador_comprobante_electronico
        string xml_sunat
        string pdf_sunat
        string hash_comprobante
        string observaciones_sunat
    }

    %% Entidad DETALLE_COMPRA
    com_DETALLE_COMPRA {
        uint64 id PK
        uint64 compra_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        decimal precio_unitario "precision(16,4)"
        decimal cantidad "precision(16,4)"
        decimal cantidad_unidad_control "precision(16,4)"
        decimal subtotal "precision(16,4)"
        decimal descuento_porcentual "precision(16,4)"
        decimal descuento_importe "precision(16,4)"
        bool recibido
        decimal cantidad_recibida "precision(16,4)"
        string notas_recepcion
        bool eliminado
        uint64 usuario_eliminacion_id FK
        string motivo_eliminacion
        string codigo_producto_sunat
        string unidad_medida_sunat
        decimal igv_linea "precision(16,4)"
        decimal isc_linea "precision(16,4)"
    }

    %% Entidad PROVEEDOR
    com_PROVEEDOR {
        uint64 id PK
        uint64 empresa_id FK
        string tipo_documento
        string numero_documento
        string nombre_razon_social
        string telefono
        string direccion
        string email
        string contacto_nombre
        string contacto_telefono
        string contacto_email
        bool activo
    }

    %% Entidad ORDEN_COMPRA
    com_ORDEN_COMPRA {
        uint64 id PK
        uint64 empresa_id FK
        string codigo
        uint64 proveedor_id FK
        time fecha_emision
        time fecha_entrega_esperada
        uint64 estado_id FK
        decimal subtotal "precision(16,4)"
        decimal impuestos "precision(16,4)"
        decimal total "precision(16,4)"
        string condiciones_entrega
        uint64 aprobador_id FK
        time fecha_aprobacion
        bool activo
    }

    %% Entidad DETALLE_ORDEN_COMPRA
    com_DETALLE_ORDEN_COMPRA {
        uint64 id PK
        uint64 orden_compra_id FK
        uint64 producto_id FK
        decimal cantidad "precision(16,4)"
        decimal precio_unitario "precision(16,4)"
        string unidad_medida
        time fecha_entrega
        decimal subtotal "precision(16,4)"
    }

    %% Entidad PRODUCTO_PROVEEDOR
    com_PRODUCTO_PROVEEDOR {
        uint64 id PK
        uint64 producto_id FK
        uint64 presentacion_id FK
        uint64 proveedor_id FK
        decimal precio_unitario "precision(16,4)"
        int tiempo_entrega_promedio
        bool es_proveedor_preferido
        decimal calificacion "precision(16,4)"
        time ultima_compra
        bool activo
        decimal cantidad_minima_pedido "precision(16,4)"
    }

    %% Entidad RECEPCION
    com_RECEPCION {
        uint64 id PK
        uint64 empresa_id FK
        uint64 compra_id FK
        uint64 estado_id FK
        time fecha_recepcion
        string numero_guia
        string transportista
        string placa_vehiculo
        string notas
        uint64 usuario_recepcion_id FK
        bool parcial
        bool finalizada
        time fecha_finalizacion
        uint64 usuario_finalizacion_id FK
    }

    %% Entidad DETALLE_RECEPCION
    com_DETALLE_RECEPCION {
        uint64 id PK
        uint64 detalle_compra_id FK
        uint64 recepcion_id FK
        decimal cantidad_recibida "precision(16,4)"
        string lote
        time fecha_vencimiento_lote
        uint64 estado_id FK
        string motivo_rechazo
        uint64 usuario_recepcion_id FK
        time fecha_recepcion
    }

    %% Entidad DETALLE_RECEPCION_LOTE
    com_DETALLE_RECEPCION_LOTE {
        uint64 id PK
        uint64 detalle_recepcion_id FK
        uint64 lote_id FK
        decimal cantidad_lote "precision(16,4)"
        string lote
        time fecha_vencimiento
        time fecha_registro
        uint64 usuario_registro_id FK
        string notas
    }

    %% Entidad IMAGEN_DOCUMENTO_COMPRA
    com_IMAGEN_DOCUMENTO_COMPRA {
        uint64 id PK
        uint64 compra_id FK
        string tipo_imagen
        string url_imagen
        time fecha_carga
        uint64 usuario_carga_id FK
        string observaciones
    }

    %% Entidad HISTORIAL_PRECIOS_PROVEEDOR
    com_HISTORIAL_PRECIOS_PROVEEDOR {
        uint64 id PK
        uint64 producto_proveedor_id FK
        decimal precio_anterior "precision(16,4)"
        decimal precio_nuevo "precision(16,4)"
        time fecha_cambio
        uint64 usuario_id FK
    }

    %% Entidades de Estado
    com_ESTADO_DOCUMENTO {
        uint64 id PK
        string codigo
        string nombre
        string descripcion
        bool activo
    }

    com_ESTADO_ORDEN {
        uint64 id PK
        string codigo
        string nombre
        string descripcion
        bool activo
    }

    com_ESTADO_RECEPCION {
        uint64 id PK
        string codigo
        string nombre
        string descripcion
        bool activo
    }

    com_ESTADO_DETALLE_RECEPCION {
        uint64 id PK
        string codigo
        string nombre
        string descripcion
        bool activo
    }

    %% Entidad TIPO_DOCUMENTO_COMPRAS
    com_TIPO_DOCUMENTO_COMPRAS {
        uint64 id PK
        uint64 empresa_id FK
        string codigo
        string nombre
        string descripcion
        bool activo
    }
::: 