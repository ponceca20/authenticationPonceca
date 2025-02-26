::: mermaid
erDiagram
    %% Leyenda:
    %% PK: Primary key, FK: Foreign key

    %% Datos Personales
    PERSONA {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        string documento_tipo "Tipo: DNI, Cédula, etc."
        string documento_numero "Número de documento"
        string foto "Foto"
        string nombre "Nombres"
        string apellidos "Apellidos"
        string email "Correo de contacto"
        string telefono "Número de contacto"
        string direccion "Dirección física"
        string ciudad "Ciudad de residencia"
        string pais "País de residencia"
        date fecha_nacimiento "Fecha de nacimiento"
    }

    %% USUARIO: cada registro incluye el sistema al que pertenece la cuenta
    USUARIO {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 persona_id FK "Referencia a PERSONA"
        uint64 sistema_id FK "Referencia al SISTEMA (para acceso con DNI y contraseña)"
        string password_hash "Contraseña específica (null si OAuth)"
        uint64 empresa_id FK "Referencia a EMPRESA (Microservicio)"
        bool activo
        uint64 creado_por FK
        uint64 actualizado_por FK
        timestamp last_login_at "Último acceso"
    }

    %% SISTEMA configurado vía formulario
    SISTEMA {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        string nombre "Nombre del sistema (configurado vía formulario)"
        string url "Ruta de acceso (configurada vía formulario; roles y módulos deben estar predefinidos)"
    }

    ROL {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 sistema_id FK "Referencia a SISTEMA para rol global"
        string codigo "Ej: ADMIN, VENDEDOR"
        string nombre
        bool activo
    }

    %% Módulos y atribución de roles en módulos
    MODULO {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 sistema_id FK "Referencia al SISTEMA"
        string codigo "Ej: VENTAS, INVENTARIO"
        string nombre
        string ruta "Ruta en la aplicación"
        bool activo
        int orden "Orden de visualización"
    }

    ROL_MODULO {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 rol_id FK "Relación con ROL"
        uint64 modulo_id FK "Relación con MODULO"
        bool acceso "Acceso asignado"
    }

    %% Entidades de seguridad y sesión
    SESION {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 usuario_id FK
        string token "JWT token"
        string refresh_token "Token de renovación"
        timestamp expires_at "Fecha de expiración"
        timestamp last_accessed_at "Último acceso durante la sesión"
        bool activa
    }

    %% USUARIO_EMPRESA: se elimina el campo sistema_id ya que USUARIO define a qué sistema pertenece
    USUARIO_EMPRESA {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 usuario_id FK "Referencia al usuario"
        uint64 empresa_id FK "Referencia a EMPRESA (Microservicio)"
        uint64 rol_id FK "Rol asignado a usuario en la empresa/sistema"
    }

    %% Relaciones
    USUARIO ||--o{ SESION : "tiene"
    SISTEMA ||--o{ MODULO : "tiene"
    ROL ||--o{ ROL_MODULO : "asigna"
    MODULO ||--o{ ROL_MODULO : "autoriza"
    USUARIO ||--o{ USUARIO_EMPRESA : "tiene acceso a"
    SISTEMA ||--o{ USUARIO : "define cuenta para"
    ROL ||--|{ USUARIO_EMPRESA : "asignado en"
:::