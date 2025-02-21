::: mermaid
erDiagram
    %% Relaciones principales
    EMPRESA ||--o{ caj_CAJA : tiene
    EMPRESA ||--o{ caj_CATEGORIA_MOVIMIENTO : tiene
    EMPRESA ||--o{ caj_NIVEL_APROBACION : tiene
    EMPRESA ||--o{ caj_SISTEMA_ORIGEN : tiene
    caj_CAJA ||--o{ caj_TURNO : "tiene"
    caj_TURNO ||--o{ caj_MOVIMIENTO_CAJA : "registra"
    caj_MOVIMIENTO_CAJA }o--|| caj_CATEGORIA_MOVIMIENTO : "pertenece a"
    caj_MOVIMIENTO_CAJA ||--o{ caj_DOCUMENTO_MOVIMIENTO_CAJA : "tiene"
    caj_TURNO ||--o| caj_CIERRE_CAJA : "finaliza con"
    caj_MOVIMIENTO_CAJA }o--|| caj_SISTEMA_ORIGEN : "proviene de"
    caj_MOVIMIENTO_CAJA }o--|| caj_ESTADO_MOVIMIENTO : "tiene"
    caj_CATEGORIA_MOVIMIENTO }o--|| caj_NIVEL_APROBACION : "requiere nivel"
    caj_MOVIMIENTO_CAJA ||--o{ caj_HISTORIAL_MOVIMIENTO : "registra"

    %% Nueva relación para canjes
    caj_MOVIMIENTO_CAJA ||--o| caj_MOVIMIENTO_CAJA : "se relaciona con"

    %% Nueva relación para cierre diario
    caj_CAJA ||--o{ caj_CIERRE_DIARIO : "tiene"
    caj_CIERRE_DIARIO ||--o{ caj_CIERRE_CAJA : "agrupa"

    caj_CAJA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre "size:100,index"
        string ubicacion "size:150"
        string codigo "size:50,unique,index"
        bool activo "default:true"
        json configuracion "Configuración general de la caja"
        string ip_address "IP del terminal"
        bool deleted
    }

    caj_TURNO {
        uint64 id PK
        uint64 caja_id FK
        uint64 usuario_id FK
        time fecha_inicio
        time fecha_fin
        string notas "size:500"
        bool requiere_arqueo
        bool deleted
    }

    caj_MOVIMIENTO_CAJA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 turno_id FK
        uint64 categoria_id FK
        uint64 estado_id FK
        uint64 sistema_origen_id FK
        uint64 usuario_autorizacion_id FK
        uint64 usuario_id FK
        uint64 documento_id FK "Referencia al documento en microservicio de pagos"
        string tipo "INGRESO|EGRESO"
        decimal monto "precision(16,4)"
        string concepto "size:255"
        string dni_beneficiario "size:20"
        string nombre_beneficiario "size:150"
        time fecha_movimiento
        json metadata
        bool deleted
        string estado_validacion "PENDIENTE|EN_VALIDACION|VALIDADO|RECHAZADO"
        uint64 metodo_pago_id FK
        uint64 movimiento_relacionado_id FK "Para canjes"
        string tipo_relacion "CANJE|DEVOLUCION|OTRO"
        string moneda "PEN|USD|EUR"
    }

    caj_CATEGORIA_MOVIMIENTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        enum tipo "INGRESO|EGRESO"
        bool requiere_aprobacion
        bool requiere_comprobante
        decimal limite_monto "precision(16,4)"
        bool activo
        uint8 orden
        uint64 nivel_aprobacion_id FK
        json validaciones
        string icono
        uint64 categoria_padre_id FK
    }

    caj_DOCUMENTO_MOVIMIENTO_CAJA {
        uint64 id PK
        uint64 movimiento_id FK
        enum tipo "FACTURA|TICKET|OTRO|IMAGEN"
        string nombre_archivo
        string url_documento
        string mime_type
        uint64 tamano_archivo
        string miniatura_url
        time fecha_documento
        bool validado
        uint64 usuario_validacion_id FK "Referencia a MS Usuarios"
    }

    caj_CIERRE_CAJA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 turno_id FK
        uint64 cierre_diario_id FK "Referencia al cierre diario"
        time fecha_cierre
        decimal saldo_inicial "precision(16,4)"
        decimal saldo_final "precision(16,4)"
        decimal total_ingresos "precision(16,4)"
        decimal total_egresos "precision(16,4)"
        decimal diferencia "precision(16,4)"
        string notas
        uint64 usuario_id FK
        bool arqueo_realizado
        uint64 supervisor_arqueo_id FK
        json desglose_medios_pago
        string tipo_cierre "PARCIAL|FINAL"
    }

    caj_NIVEL_APROBACION {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        decimal limite_monto "precision(16,4)"
        bool requiere_supervisor
        uint8 nivel_prioridad
        bool activo
    }

    caj_HISTORIAL_MOVIMIENTO {
        uint64 id PK
        uint64 movimiento_id FK
        enum estado "PENDIENTE|APROBADO|RECHAZADO"
        string comentario
        time fecha_cambio
        uint64 usuario_id FK "Referencia a MS Usuarios"
    }

    caj_SISTEMA_ORIGEN {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string codigo
        bool activo
        json configuracion_integracion
    }

    caj_ESTADO_MOVIMIENTO {
        uint64 id PK
        enum nombre "BORRADOR|PENDIENTE|APROBADO|RECHAZADO|ANULADO"
        string descripcion
        bool permite_edicion
        bool requiere_comentario
    }

    caj_CIERRE_DIARIO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 caja_id FK
        date fecha_cierre "Fecha del día del cierre"
        time hora_cierre
        decimal saldo_inicial_dia "precision(16,4)"
        decimal saldo_final_dia "precision(16,4)"
        decimal total_ingresos_dia "precision(16,4)"
        decimal total_egresos_dia "precision(16,4)"
        decimal diferencia_dia "precision(16,4)"
        string estado "PENDIENTE|EN_PROCESO|CERRADO|ANULADO"
        uint64 usuario_cierre_id FK
        uint64 supervisor_id FK
        string notas "size:500"
        json resumen_medios_pago
        bool cuadre_validado
        bool deleted
    }

    caj_CAJA ||--o{ caj_CAJA_IMPRESORA : configura

    caj_CAJA_IMPRESORA {
        uint64 id PK
        uint64 caja_id FK
        string nombre "Ejemplo: Impresora Tickets"
        string tipo "TICKET/COMANDA/FACTURA"
        string modelo
        string ip_address
        string puerto "USB/RED/SERIE"
        json configuracion "Configuración específica"
        bool es_principal
        bool activo
    }
:::
