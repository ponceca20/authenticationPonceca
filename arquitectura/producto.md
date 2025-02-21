::: mermaid
erDiagram
    EMPRESA ||--o{ prod_PRODUCTO : tiene
    EMPRESA ||--o{ prod_CATEGORIA : tiene
    EMPRESA ||--o{ prod_TIPO_PRODUCTO : tiene
    EMPRESA ||--o{ prod_AREA_DESPACHO : tiene
    EMPRESA ||--o{ prod_UNIDAD : tiene
    EMPRESA ||--o{ prod_CARTA : tiene
    EMPRESA ||--o{ prod_IMPUESTO : tiene
    EMPRESA ||--o{ prod_GRUPO_CONTEO : tiene
    prod_PRODUCTO ||--o{ prod_PRESENTACION : tiene
    prod_PRODUCTO }o--|| prod_CATEGORIA : pertenece
    prod_PRODUCTO }o--|| prod_TIPO_PRODUCTO : clasifica
    prod_PRODUCTO }o--|| prod_AREA_DESPACHO : despacha
    prod_PRODUCTO }o--|| prod_UNIDAD : controla
    prod_PRODUCTO ||--o{ prod_COMENTARIO_PRODUCTO : permite
    prod_PRODUCTO ||--o{ prod_PRODUCTO_GRUPO_CONTEO : "puede tener"
    prod_PRODUCTO ||--o{ prod_PRODUCTO_MEDIA : "tiene"
    prod_PRESENTACION ||--o{ prod_PRESENTACION_MEDIA : "tiene"
    prod_CARTA ||--|{ prod_FAMILIA : "contiene"
    prod_FAMILIA ||--|{ prod_PRODUCTO_FAMILIA : "contiene"
    prod_PRODUCTO ||--|{ prod_PRODUCTO_FAMILIA : "está en"
    prod_PRODUCTO }o--|| prod_IMPUESTO : "puede tener"

    prod_PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"        
        string nombre
        string descripcion
        string estado "ACTIVO|INACTIVO|ELIMINADO"
        uint64 tipo_producto_id FK
        decimal precio_base "precision(16,4)"
        uint64 categoria_id FK
        uint64 area_despacho_id FK
        uint64 unidad_control_id FK
        int tiempo_preparacion
        string codigo UK
        string codigo_barras "Código de barras del producto"
        string tipo_gestion "TPV|ALMACEN|AMBOS|NINGUNO"
        string tipo_control "STOCK|STOCK_TEMPORAL|SERIE|VENCIMIENTO|LOTE"
        string politica_lotes "FIFO|FEFO|LIFO"
        bool es_transformable
        bool requiere_receta
        uint64 impuesto_id FK
    }

    prod_PRESENTACION {
        uint64 id PK
        string nombre       
        decimal cantidad_unidad_control "precision(16,4)"
        decimal precio "precision(16,4)"
        uint64 producto_id FK    
        string estado "ACTIVO|INACTIVO|ELIMINADO"
        bool es_presentacion_base
        string codigo_barras "Código de barras de la presentación"
    }

    prod_PRODUCTO_MEDIA {
        uint64 id PK
        uint64 producto_id FK
        string url
        string tipo "IMAGEN|VIDEO"
        int orden
        string estado "ACTIVO|INACTIVO"
        bool es_principal
    }

    prod_PRESENTACION_MEDIA {
        uint64 id PK
        uint64 presentacion_id FK
        string url
        string tipo "IMAGEN|VIDEO"
        int orden
        string estado "ACTIVO|INACTIVO"
        bool es_principal
    }

    prod_TIPO_PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string estado "ACTIVO|INACTIVO"
        string descripcion "Ayuda al usuario a entender el tipo"
    }

    prod_AREA_DESPACHO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre        
        string estado "ACTIVO|INACTIVO"
        string impresora_asociada 
        string tipo_impresion "IMPRESORA|PANTALLA"
    }

    prod_CATEGORIA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre      
        string estado "ACTIVO|INACTIVO"
    }

    prod_UNIDAD {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string simbolo      
        bool es_decimal
        string estado "ACTIVO|INACTIVO"
    }

    prod_COMENTARIO_PRODUCTO {
        uint64 id PK
        string comentario
        string estado "ACTIVO|INACTIVO"
        uint64 producto_id FK "Null si es genérico"
        string tipo "GENERICO|ESPECIFICO"
    }

    prod_GRUPO_CONTEO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        string estado "ACTIVO|INACTIVO"
        uint64 almacen_id FK "Referencia al ID del almacén en el microservicio de inventario"
        int orden_conteo "Orden para la secuencia de conteo"
    }

    prod_CARTA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre 
        string descripcion
        string estado "ACTIVO|INACTIVO"
        string imagen
        int orden_visualizacion
        string control_horario "CON_HORARIO|SIN_HORARIO"      
        time hora_inicio_servicio
        time hora_fin_servicio
    }
    
    prod_FAMILIA {
        uint64 id PK
        uint64 carta_id FK
        string nombre        
        int orden_visualizacion
        string estado "ACTIVO|INACTIVO"
        string imagen
    }
    
    prod_PRODUCTO_FAMILIA {
        uint64 id PK
        uint64 producto_id FK
        uint64 familia_id FK      
        int orden_visualizacion
        bool destacado       
        string estado "DISPONIBLE|NO_DISPONIBLE"
        string control_horario "CON_HORARIO|SIN_HORARIO"
        time hora_inicio_disponibilidad
        time hora_fin_disponibilidad
    }

    prod_IMPUESTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        decimal porcentaje "precision(16,4)"
        string estado "ACTIVO|INACTIVO|ELIMINADO"
    }

    prod_PRODUCTO_GRUPO_CONTEO {
        uint64 id PK
        uint64 producto_id FK
        uint64 grupo_conteo_id FK
        string estado "ACTIVO|INACTIVO"
    }

    prod_GRUPO_CONTEO ||--o{ prod_PRODUCTO_GRUPO_CONTEO : "contiene"
::: 