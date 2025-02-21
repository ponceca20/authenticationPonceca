::: mermaid
erDiagram
    EMPRESA ||--o{ usr_USUARIO : tiene
    EMPRESA ||--o{ usr_ROL : tiene
    usr_USUARIO ||--o{ usr_SESION : tiene
    usr_USUARIO ||--o{ usr_USUARIO_ROL : tiene
    usr_USUARIO_ROL }o--|| usr_ROL : asignado
    usr_ROL ||--o{ usr_ROL_PERMISO : tiene
    usr_ROL_PERMISO }o--|| usr_PERMISO : asignado
    usr_PERMISO ||--o{ usr_MODULO_PERMISO : tiene
    usr_MODULO_PERMISO }o--|| usr_MODULO : pertenece

    usr_USUARIO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string email "Email único del usuario"
        string password_hash "Hash de contraseña (null si usa Google)"
        string dni
        string nombre
        string apellidos
        string telefono
        bool activo
        string google_id "ID único de Google (para OAuth)"
        string foto_url
        time ultimo_acceso
        time ultimo_cambio_password "Fecha del último cambio de contraseña"
        string zona_horaria "Zona horaria del usuario"
        uint64 creado_por FK
        uint64 actualizado_por FK
    }

    usr_SESION {
        uint64 id PK
        uint64 usuario_id FK
        string token "JWT token"
        string refresh_token
        string dispositivo "Info del dispositivo/navegador"
        string ip_address
        string zona_horaria_login "Zona horaria desde donde se inició sesión"
        time fecha_creacion
        time fecha_expiracion
        bool activa
    }

    usr_ROL {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string codigo "ADMIN|VENDEDOR|CAJERO|etc"
        string nombre
        string descripcion
        bool activo
    }

    usr_USUARIO_ROL {
        uint64 id PK
        uint64 usuario_id FK
        uint64 rol_id FK
        time fecha_asignacion
        bool activo
        uint64 asignado_por FK
    }

    usr_PERMISO {
        uint64 id PK
        string codigo "CREATE|READ|UPDATE|DELETE"
        string nombre
        string descripcion
        bool activo
    }

    usr_ROL_PERMISO {
        uint64 id PK
        uint64 rol_id FK
        uint64 permiso_id FK
        uint64 modulo_id FK
        bool activo
    }

    usr_MODULO {
        uint64 id PK
        string codigo "VENTAS|INVENTARIO|REPORTES|etc"
        string nombre
        string descripcion
        string ruta "Ruta en la aplicación"
        bool activo
        uint64 modulo_padre_id FK "Para submodulos"
        int orden "Orden de visualización"
    }

    usr_MODULO_PERMISO {
        uint64 id PK
        uint64 modulo_id FK
        uint64 permiso_id FK
        bool activo
    }
::: 