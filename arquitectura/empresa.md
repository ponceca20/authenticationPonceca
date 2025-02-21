::: mermaid
erDiagram
    emp_EMPRESA ||--o{ emp_CONFIGURACION : configura
    emp_EMPRESA ||--o{ emp_EMPRESA_RELACION : tiene
    emp_EMPRESA ||--o{ emp_SUSC_SUSCRIPCION : tiene
    emp_SUSC_SUSCRIPCION ||--o{ emp_SUSC_PAGO : registra
    emp_SUSC_PLAN ||--o{ emp_SUSC_SUSCRIPCION : define
    emp_SUSC_SUSCRIPCION ||--o{ emp_SUSC_DOCUMENTO : genera

    emp_EMPRESA {
        uint64 id PK
        string codigo UK "Código único de empresa"
        string ruc UK "RUC de la empresa"
        string razon_social
        string nombre_comercial
        string direccion_fiscal
        string telefono
        string email
        string zona_horaria "Zona horaria de la empresa"
        json datos_sunat "Usuario, clave, certificado, etc"
        json datos_contacto "Web, redes sociales, etc"
        json logos "URLs de los diferentes logos"
        float64 latitud "Ubicación geográfica"
        float64 longitud "Ubicación geográfica"
        json horario "Horarios de atención"
        string estado "ACTIVO|INACTIVO|ELIMINADO"
    }

    emp_EMPRESA_RELACION {
        uint64 id PK
        uint64 empresa_matriz_id FK "Empresa principal"
        uint64 empresa_sucursal_id FK "Empresa relacionada"
        string tipo_relacion "SUCURSAL|FRANQUICIA|ASOCIADA"  
        string estado "ACTIVO|INACTIVO"
    }

    emp_CONFIGURACION {
        uint64 id PK
        uint64 empresa_id FK
        string clave UK "Identificador de configuración"
        string valor "Valor de la configuración"
        bool es_secreto "Para datos sensibles"
    }

    emp_SUSC_PLAN {
        uint64 id PK
        string nombre
        string descripcion
        decimal precio_mensual "precision(16,4)"
        decimal precio_anual "precision(16,4)"
        bool es_free
        bool activo
    }

    emp_SUSC_SUSCRIPCION {
        uint64 id PK
        uint64 empresa_id FK
        uint64 plan_id FK
        date fecha_inicio
        date fecha_fin
        string periodicidad "FREE|MENSUAL|ANUAL"
        string estado "ACTIVA|PENDIENTE|SUSPENDIDA|CANCELADA"
        bool notificado
        int limite_usuarios "Número máximo de usuarios permitidos"
    }

    emp_SUSC_PAGO {
        uint64 id PK
        uint64 suscripcion_id FK
        decimal monto "precision(16,4)"
        date fecha_pago
        date periodo_inicio
        date periodo_fin
        string estado "PENDIENTE|PAGADO|CANCELADO"
        string metodo_pago "EFECTIVO|TARJETA|TRANSFERENCIA|DEPOSITO"
        string referencia_pago
    }

    emp_SUSC_DOCUMENTO {
        uint64 id PK
        uint64 empresa_id FK
        uint64 suscripcion_id FK "Referencia a la suscripción"
        string tipo "FACTURA|BOLETA|NOTA_CREDITO|NOTA_DEBITO"
        decimal monto_total "precision(16,4)"
        decimal monto_pendiente "precision(16,4)"
        time fecha_emision
        time fecha_vencimiento
        bool es_credito
        string estado_pago "PENDIENTE|PAGADO|VENCIDO|ANULADO"
    }
::: 