::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ ven_SALON : tiene
    EMPRESA ||--o{ ven_VENTA : tiene
    EMPRESA ||--o{ ven_CLIENTE : tiene
    EMPRESA ||--o{ ven_SECUENCIA_DOCUMENTO : tiene
    ven_SALON ||--o{ ven_MESA : contiene
    ven_VENTA ||--o{ ven_DETALLE_VENTA : contiene
    ven_VENTA }o--|| ven_CLIENTE : realiza
    ven_DETALLE_VENTA ||--o{ ven_DETALLE_VENTA_MODIFICADOR : "puede tener"
    ven_VENTA ||--o{ ven_VENTA : "dividida en"
    ven_VENTA }|--|| ven_TIPO_DOCUMENTO_VENTA : "tiene"
    ven_TIPO_DOCUMENTO_VENTA ||--o{ ven_SECUENCIA_DOCUMENTO : "tiene"
    ven_VENTA }|--|| ven_TIPO_OPERACION_SUNAT : tiene
    ven_VENTA }|--|| ven_MONEDA : tiene
    ven_DETALLE_VENTA ||--o{ ven_DETALLE_VENTA_LOTE : "utiliza"

    %% Entidades base (sin cambios)
    ven_SALON {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        string estado "ACTIVO|INACTIVO"
        uint8 orden
        string codigo_qr
    }

    ven_MESA {
        uint64 id PK
        uint64 salon_id FK
        decimal numero_mesa "precision(16,4)"
        string descripcion
        uint8 capacidad
        string estado "ACTIVO|INACTIVO"
        string estado_ocupacion "LIBRE|OCUPADA|RESERVADA|REQUIERE_LIMPIEZA"
        string codigo_qr
    }

    %% Entidad VENTA actualizada para SUNAT
    ven_VENTA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 cliente_id FK
        uint64 venta_origen_id FK
        uint64 caja_id FK
        uint64 turno_id FK
        decimal numero_mesa "precision(16,4)"
        uint64 tipo_documento_venta_id FK
        string serie
        string numero
        string serie_correlativo "FXXX-XXXXXXXX"
        decimal subtotal "precision(16,4)"
        decimal impuestos "precision(16,4)"
        decimal total "precision(16,4)"
        string estado "PENDIENTE|EN_PROCESO|COMPLETADA|ANULADA|CANCELADA"
        time fecha_emision
        time fecha_vencimiento
        string notas
        decimal cargo_delivery "precision(16,4)"
        string codigo_seguimiento
        uint64 usuario_creacion_id FK
        decimal descuento_porcentual "precision(16,4)"
        decimal descuento_importe "precision(16,4)"
        bool emplatado
        bool facturado
        bool enviado_sunat
        string estado_sunat "ACEPTADO|RECHAZADO|PENDIENTE"
        string hash_cpe
        string cdr_sunat
        string codigo_qr "QR según spec SUNAT"
        string tipo_operacion "0101=Venta interna"
        string moneda "PEN|USD"
        decimal tipo_cambio "precision(16,4)"
        string documento_relacionado "Para NC/ND"
        string codigo_motivo "Para NC/ND"
        string descripcion_motivo "Para NC/ND"
        bool gratuita
        string metodo_pago "CONTADO|CREDITO"
        string condiciones_pago
        bool deleted
        string tipo_division "PORCENTAJE|MONTO|ITEMS"
        uint8 numero_division 
        uint8 total_divisiones
        string estado_division "PENDIENTE|EN_PROCESO|COMPLETADA"
        decimal porcentaje_division "precision(16,4)"
        uint64 usuario_division_id FK
        time fecha_division
        string notas_division
        bool requiere_credito
        uint64 documento_credito_id FK
        uint64 recepcion_oc_id FK "Referencia a la OC origen"
        string estado_origen "GENERADO_POR_OC|MANUAL"
    }

    %% DETALLE_VENTA actualizado para SUNAT
    ven_DETALLE_VENTA {
        uint64 id PK
        uint64 venta_id FK
        uint64 presentacion_id FK
        uint64 oferta_id FK
        uint64 oferta_presentacion_id FK
        bool es_parte_oferta
        uint8 posicion_en_oferta
        string nombre_item
        string codigo_producto_sunat
        string unidad_medida_sunat
        decimal cantidad_solicitada "precision(16,4)"
        decimal cantidad_unidad_base "precision(16,4)"
        decimal precio_unitario_normal "precision(16,4)"
        decimal precio_unitario_oferta "precision(16,4)"
        decimal subtotal_linea "precision(16,4)"
        string tipo_afectacion_igv "GRAVADO|EXONERADO|INAFECTO"
        decimal igv_linea "precision(16,4)"
        decimal isc_linea "precision(16,4)"
        decimal porcentaje_descuento "precision(16,4)"
        decimal importe_descuento "precision(16,4)"
        bool gratuita
        uint8 secuencia_preparacion
        bool marcado_urgente
        string estado_preparacion "PENDIENTE|EN_PROCESO|PREPARADO|ENTREGADO|CANCELADO"
        uint8 nivel_prioridad
        time hora_entrega_prometida
        bool eliminado_logico
        time fecha_eliminacion
        uint64 usuario_eliminacion_id FK
        string razon_eliminacion
        decimal cantidad_dividida "precision(16,4)"
        uint64 detalle_venta_origen_id FK
        decimal porcentaje_division_item "precision(16,4)"
        uint64 usuario_modificacion_id FK
        time fecha_modificacion
        string razon_modificacion
    }

    %% CLIENTE actualizado para SUNAT
    ven_CLIENTE {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string tipo_documento_sunat "6=RUC, 1=DNI"
        string numero_documento
        string nombre_razon_social
        string telefono
        string direccion
        string email
        string ubigeo
        string codigo_pais "PE"
        bool acepta_marketing
        time ultima_visita
    }

    %% Entidades de configuración SUNAT
    ven_TIPO_DOCUMENTO_VENTA {
        uint64 id PK
        string codigo "01=Factura, 03=Boleta, 07=NC, 08=ND"
        string nombre
        string serie_formato
        string numero_formato
        uint64 longitud_numero
        string prefijo_serie
        bool requiere_cliente
        bool activo
        string descripcion
        string[] codigos_motivo
        string[] descripciones_motivo
    }

    ven_TIPO_OPERACION_SUNAT {
        string codigo "0101, 0102, etc."
        string descripcion
        bool activo
    }

    ven_MONEDA {
        string codigo "PEN/USD"
        string descripcion
        bool activo
    }

    %% Entidades sin cambios
    ven_SECUENCIA_DOCUMENTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 tipo_documento_id FK
        string serie
        uint64 ultimo_numero
        uint64 establecimiento_id FK
        bool activo
        time fecha_ultima_actualizacion
        uint64 usuario_ultima_actualizacion_id FK
    }

    ven_DETALLE_VENTA_MODIFICADOR {
        uint64 id PK
        uint64 detalle_venta_id FK
        uint64 modificador_id FK
        uint64 modificador_opcion_id FK
        decimal precio_adicional "precision(16,4)"
    }

    %% Nueva entidad para relacionar detalles de venta con lotes
    ven_DETALLE_VENTA_LOTE {
        uint64 id PK
        uint64 detalle_venta_id FK
        uint64 lote_id FK
        decimal cantidad_lote "precision(16,4)"
        time fecha_registro
        uint64 usuario_registro_id FK
    }
:::