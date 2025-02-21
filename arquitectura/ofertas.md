::: mermaid
erDiagram
    EMPRESA ||--o{ off_OFERTA : tiene
    EMPRESA ||--o{ off_MODIFICADOR : tiene
    
    off_OFERTA ||--|{ off_GRUPO_OFERTA : "contiene"
    off_GRUPO_OFERTA ||--|{ off_OFERTA_PRESENTACION : "organiza"
    off_OFERTA_PRESENTACION }|--|| PRESENTACION : "usa" 
    PRESENTACION ||--o{ off_PRESENTACION_MODIFICADOR : "tiene"
    off_PRESENTACION_MODIFICADOR ||--o{ off_MODIFICADOR : "usa"
    off_MODIFICADOR ||--|{ off_MODIFICADOR_OPCION : "tiene"   
    FAMILIA ||--|{ off_OFERTA_FAMILIA : "contiene"
    off_OFERTA ||--|{ off_OFERTA_FAMILIA : "está en"

    off_OFERTA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 establecimiento_id FK
        string tipo_categoria
        string nombre
        string descripcion
        string imagen_url
        decimal precio_base "precision(16,4)"
        time fecha_inicio
        time fecha_fin
        bool tiene_control_fechas
        bool activo
        bool deleted
    }
    
    off_GRUPO_OFERTA {
        uint64 id PK
        uint64 oferta_id FK
        string nombre
        string descripcion
        uint8 min_selecciones
        uint8 max_selecciones
        decimal precio_base "precision(16,4)"
        int orden_visualizacion
        bool activo
        bool deleted
    }
    
    off_OFERTA_PRESENTACION {
        uint64 id PK
        uint64 oferta_id FK
        uint64 grupo_oferta_id FK
        uint64 presentacion_id FK
        decimal precio_adicional "precision(16,4)"
        decimal stock_diario "precision(16,4)"
        decimal stock_actual "precision(16,4)"
        bool activo
        bool deleted
    }

    off_PRESENTACION_MODIFICADOR {
        uint64 id PK
        uint64 presentacion_id FK
        uint64 modificador_id FK
        bool required
        int orden_visualizacion
        bool activo
        bool deleted
    }

    off_MODIFICADOR {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        string tipo
        uint8 min_selecciones
        uint8 max_selecciones
        bool multiple_seleccion
        bool activo
        bool deleted
    }

    off_MODIFICADOR_OPCION {
        uint64 id PK
        uint64 modificador_id FK
        string nombre
        string descripcion
        decimal precio_adicional "precision(16,4)"
        int orden_visualizacion
        bool activo
        bool deleted
        string imagen_url 
    }

    off_OFERTA_FAMILIA {
        uint64 id PK
        uint64 oferta_id FK
        uint64 familia_id FK
        int orden_visualizacion
        bool destacado
        bool disponible
        bool tiene_control_horario
        time hora_inicio_disponibilidad
        time hora_fin_disponibilidad
        bool activo
        bool deleted
    }
::: 