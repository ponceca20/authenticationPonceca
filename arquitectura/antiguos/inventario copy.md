::: mermaid
erDiagram
    EMPRESA ||--o{ inv_TRANSACCION_INVENTARIO : tiene
    EMPRESA ||--o{ inv_TIPO_MOVIMIENTO : tiene
    EMPRESA ||--o{ inv_ALMACEN : tiene
    EMPRESA ||--o{ inv_KARDEX : tiene
    EMPRESA ||--o{ inv_CONTEO_INVENTARIO : tiene
    inv_TRANSACCION_INVENTARIO ||--o{ inv_DETALLE_TRANSACCION : contiene
    inv_TRANSACCION_INVENTARIO }o--|| inv_TIPO_MOVIMIENTO : clasifica
    inv_TRANSACCION_INVENTARIO }o--|| inv_ALMACEN : afecta
    inv_MOVIMIENTO_ALMACEN }o--|| inv_ALMACEN : registra
    inv_MOVIMIENTO_ALMACEN }o--|| inv_KARDEX : afecta
    inv_CONTEO_INVENTARIO ||--o{ inv_DETALLE_CONTEO : contiene
    inv_CONTEO_INVENTARIO }o--|| inv_ALMACEN : realiza_en

    inv_TRANSACCION_INVENTARIO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 tipo_movimiento_id FK
        uint64 almacen_origen_id FK
        uint64 almacen_destino_id FK
        string numero_documento
        time fecha_movimiento
        string estado
        string motivo
        string notas
        uint64 usuario_creacion_id FK
        uint64 usuario_autorizacion_id FK
        uint64 usuario_recepcion_id FK
        time fecha_recepcion
        string origen_transaccion "Ventas/Compras/Producción/etc"
        string documento_origen "Referencia al documento origen"
        bool deleted
    }

    inv_MOVIMIENTO_ALMACEN {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 kardex_id FK
        uint64 almacen_id FK
        uint64 producto_id FK
        uint64 presentacion_id FK
        string tipo_movimiento "ENTRADA|SALIDA"
        decimal cantidad "precision(16,4)"
        decimal stock_anterior "precision(16,4)"
        decimal stock_actual "precision(16,4)"
        decimal costo_unitario "precision(16,4)"
        decimal costo_total "precision(16,4)"
        time fecha_movimiento
        string origen_tipo "VENTA|COMPRA|PRODUCCION|AJUSTE"
        uint64 origen_id "ID del documento origen"
        uint64 origen_id_detalle "ID del detalle del documento origen"
        string notas
        uint64 usuario_id FK
    }

    inv_DETALLE_TRANSACCION {
        uint64 id PK
        uint64 transaccion_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        string nombre_producto
        decimal cantidad_enviada "precision(16,4)"
        decimal cantidad_recibida "precision(16,4)"
        decimal cantidad_unidad_control "precision(16,4)"
        decimal costo_unitario "precision(16,4)"
        string lote
        time fecha_vencimiento
        string ubicacion_origen
        string ubicacion_destino
        string notas
        bool deleted
    }

    inv_TIPO_MOVIMIENTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string codigo
        string nombre "Entrada, Salida, Traslado, Ajuste"
        string descripcion
        bool afecta_costo
        bool requiere_autorizacion
        bool requiere_recepcion "True para traslados"
        bool activo
    }

    inv_KARDEX {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 presentacion_id FK "Referencia a Productos.PRESENTACION"
        uint64 producto_id FK "Referencia a Productos.PRODUCTO"
        uint64 almacen_id FK
        decimal stock_actual "precision(16,4)"
        decimal stock_minimo "precision(16,4)"
        decimal stock_maximo "precision(16,4)"
        decimal costo_promedio "precision(16,4)"
        decimal ultima_entrada "precision(16,4)"
        decimal ultima_salida "precision(16,4)"
        time ultimo_movimiento
        string ubicacion_preferida
        bool requiere_lote
        bool control_vencimiento
    }

    inv_ALMACEN {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string codigo
        string nombre
        string descripcion
        bool es_produccion
        bool es_despacho
        bool activo
        string ubicacion
        uint64 responsable_id FK
    }

    inv_CONTEO_INVENTARIO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 almacen_id FK "Almacén donde se realiza el conteo"
        string codigo
        time fecha_inicio
        time fecha_fin
        string estado "Programado, En Proceso, Finalizado"
        string tipo "Parcial, Total"
        string notas
        uint64 usuario_creacion_id FK
        uint64 usuario_cierre_id FK
        bool deleted
    }

    inv_DETALLE_CONTEO {
        uint64 id PK
        uint64 conteo_id FK
        uint64 presentacion_id FK
        uint64 producto_id FK
        string nombre_producto
        uint64 almacen_id FK
        decimal stock_sistema "precision(16,4)"
        decimal stock_contado "precision(16,4)"
        decimal diferencia "precision(16,4)"
        string notas
        bool ajustado
        time fecha_ajuste
        uint64 usuario_ajuste_id FK
    }

    inv_DETALLE_TRANSACCION ||--o{ inv_DETALLE_TRANSACCION_LOTE : "gestiona"
    inv_MOVIMIENTO_ALMACEN ||--o{ inv_MOVIMIENTO_ALMACEN_LOTE : "utiliza"

    inv_DETALLE_TRANSACCION_LOTE {
        uint64 id PK
        uint64 detalle_transaccion_id FK
        uint64 lote_id FK "Referencia al microservicio de lotes"
        decimal cantidad_lote "precision(16,4)"
        time fecha_registro
        uint64 usuario_registro_id FK
        string notas
    }

    inv_MOVIMIENTO_ALMACEN_LOTE {
        uint64 id PK
        uint64 movimiento_almacen_id FK
        uint64 lote_id FK "Referencia al microservicio de lotes"
        decimal cantidad_lote "precision(16,4)"
        time fecha_registro
        uint64 usuario_registro_id FK
    }
:::
