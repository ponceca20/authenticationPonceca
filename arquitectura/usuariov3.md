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

    %% USUARIO: cuenta para acceso al sistema
    USUARIO {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 persona_id FK "Referencia a PERSONA"
        string password_hash "Contraseña encriptada"
        bool activo "Estado del usuario"
        uint64 creado_por FK "Usuario que creó el registro"
        uint64 actualizado_por FK "Usuario que actualizó el registro"
    }

    %% ROL: define los permisos en el sistema
    ROL {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 empresa_id FK "Referencia a EMPRESA"
        string codigo "Ej: ADMIN, VENDEDOR"
        string nombre "Nombre del rol"
        bool activo "Estado del rol"
    }

    %% Módulos y atribución de roles en módulos
    MODULO {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        string codigo "Ej: VENTAS, INVENTARIO"
        string nombre "Nombre del módulo"
        string ruta "Ruta en la aplicación"
        bool activo "Estado del módulo"
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
        uint64 usuario_id FK "Referencia a USUARIO"
        string token "JWT token"
        string refresh_token "Token de renovación"
        timestamp fecha_expiracion "Fecha de expiración"
        bool activa "Estado de la sesión"
        string ip "Dirección IP"
    }

    %% USUARIO_EMPRESA: relación entre usuarios y empresas con roles
    USUARIO_EMPRESA {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        uint64 usuario_id FK "Referencia al usuario"
        uint64 empresa_id FK "Referencia a EMPRESA"
        uint64 rol_id FK "Rol asignado a usuario en la empresa"
        timestamp fecha_asignacion "Fecha en que se asignó el rol"
    }

    %% Control de intentos de inicio de sesión
    LOGIN_ATTEMPT {
        uint64 id PK "Identificador único"
        timestamp created_at "Fecha de creación"
        timestamp updated_at "Fecha de actualización"
        timestamp deleted_at "Fecha de borrado"
        string identifier "Email o nombre de usuario (puede no corresponder a usuario válido)"
        string ip "Dirección IP del intento"
        bool success "Si el intento fue exitoso"
    }

    %% Nota: LOGIN_ATTEMPT no tiene relaciones directas ya que registra todos los intentos 
    %% de inicio de sesión, incluyendo aquellos con credenciales inválidas o usuarios que no existen.
    %% Su propósito es detectar intentos de acceso no autorizados, ataques de fuerza bruta
    %% y proporcionar información para auditorías de seguridad del sistema.

    %% Relaciones
    PERSONA ||--o{ USUARIO : "tiene"
    USUARIO ||--o{ SESION : "tiene"
    ROL ||--o{ ROL_MODULO : "asigna"
    MODULO ||--o{ ROL_MODULO : "autoriza"
    USUARIO ||--o{ USUARIO_EMPRESA : "tiene acceso a"
    ROL ||--|{ USUARIO_EMPRESA : "asignado en"
:::