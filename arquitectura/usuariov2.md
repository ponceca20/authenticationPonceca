::: mermaid
erDiagram
  
    USUARIO ||--o{ USUARIO_EMPRESA : "tiene acceso a"
    ROL ||--|{ USUARIO_EMPRESA : "asignado en"
    
    PERSONA {
        uint64 id PK "Identificador único"
        string documento_tipo "Tipo: DNI, Cédula, etc."
        string documento_numero "Número de documento"
        string foto "foto"
        string nombre "Nombres"
        string apellidos "Apellidos"
        string email "Correo de contacto"
        string telefono "Número de contacto"
        string direccion "Dirección física"
        string ciudad "Ciudad de residencia"
        string pais "País de residencia"
        date fecha_nacimiento "Fecha de nacimiento"
        time deleted_at "Borrado lógico"
    }
    
    USUARIO {
        uint64 id PK "Identificador único global"
        uint64 persona_id FK "Referencia a PERSONA"
        string password_hash "Hash de contraseña (null si OAuth)"
        bool activo
        uint64 creado_por FK
        uint64 actualizado_por FK
        time deleted_at "Borrado lógico"
    }
    
    USUARIO_EMPRESA {
        uint64 id PK
        uint64 usuario_id FK "Referencia al usuario"
        uint64 empresa_id FK "Referencia a EMPRESA (Microservicio)"
        uint64 rol_id FK "Rol asignado en la empresa"
        time fecha_asignacion
        time deleted_at "Borrado lógico"
    }
    
    ROL {
        uint64 id PK
        uint64 empresa_id FK "Referencia a EMPRESA (Microservicio)"
        string codigo "Ej: ADMIN, VENDEDOR"
        string nombre
        bool activo
     
        time deleted_at "Borrado lógico"
    }
    
    SESION {
        uint64 id PK
        uint64 usuario_id FK
        string token "JWT token"
        string refresh_token
        time fecha_creacion
        time fecha_expiracion       
        bool activa
    }
    
    MODULO {
        uint64 id PK
        string codigo "Ej: VENTAS, INVENTARIO"
        string nombre
        string ruta "Ruta en la aplicación"
        bool activo
        int orden "Orden de visualización"
        time deleted_at "Borrado lógico"
    }
    
    ROL_MODULO {
        uint64 id PK
        uint64 rol_id FK "Relación con ROL"
        uint64 modulo_id FK "Relación con MODULO"
        bool acceso "Acceso asignado"
        time fecha_creacion
        time fecha_actualizacion
        time deleted_at "Borrado lógico"
    }
    
    USUARIO ||--o{ SESION : "tiene"
    ROL ||--o{ ROL_MODULO : "asigna"
    MODULO ||--o{ ROL_MODULO : "autoriza"
:::