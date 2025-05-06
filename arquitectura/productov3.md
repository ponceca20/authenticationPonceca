::: mermaid
erDiagram
    EMPRESA ||--o{ PRODUCTO : tiene
    EMPRESA ||--o{ CATEGORIA : tiene
    EMPRESA ||--o{ TIPO_PRODUCTO : tiene
    EMPRESA ||--o{ AREA_DESPACHO : tiene
    EMPRESA ||--o{ UNIDAD : tiene
    EMPRESA ||--o{ CARTA : tiene
    EMPRESA ||--o{ IMPUESTO : tiene
    EMPRESA ||--o{ GRUPO_CONTEO : tiene
    EMPRESA ||--o{ MONITOR_COCINA : tiene
    EMPRESA ||--o{ MONITOR_AUXILIAR : tiene
    EMPRESA ||--o{ OFERTA : tiene
    EMPRESA ||--o{ MODIFICADOR : tiene
    
    PRODUCTO ||--o{ PRESENTACION : tiene
    PRODUCTO }o--|| CATEGORIA : pertenece
    PRODUCTO }o--|| TIPO_PRODUCTO : clasifica
    PRODUCTO }o--|| AREA_DESPACHO : despacha
    PRODUCTO }o--|| UNIDAD : controla
    PRODUCTO ||--o{ COMENTARIO_PRODUCTO : permite
    PRODUCTO ||--o{ PRODUCTO_GRUPO_CONTEO : "puede tener"
    PRODUCTO ||--o{ PRODUCTO_MEDIA : "tiene"
    PRESENTACION ||--o{ PRESENTACION_MEDIA : "tiene"
    CARTA ||--|{ FAMILIA : "contiene"
    FAMILIA ||--|{ PRODUCTO_FAMILIA : "contiene"
    PRODUCTO ||--|{ PRODUCTO_FAMILIA : "está en"
    PRODUCTO }o--|| IMPUESTO : "puede tener"
    PRODUCTO }o--|| MONITOR_COCINA : "usa"
    PRODUCTO }o--|| MONITOR_AUXILIAR : "usa"
    OFERTA ||--|{ GRUPO_OFERTA : "contiene"
    GRUPO_OFERTA ||--|{ OFERTA_PRESENTACION : "organiza"
    OFERTA_PRESENTACION }|--|| PRESENTACION : "usa" 
    PRESENTACION ||--o{ PRESENTACION_MODIFICADOR : "tiene"
    PRESENTACION_MODIFICADOR ||--o{ MODIFICADOR : "usa"
    MODIFICADOR ||--|{ MODIFICADOR_OPCION : "tiene"
    PRESENTACION |o--|| OFERTA : "puede ser"
    PRESENTACION ||--o{ PRESENTACION_PESO : "puede tener"
    PRESENTACION ||--o{ PRESENTACION_TIEMPO : "puede tener"
    PRESENTACION ||--o{ PRECIO_HISTORICO : "registra cambios"
    PRESENTACION ||--o{ PRESENTACION_COMPUESTA : "contiene"
    PRESENTACION_COMPUESTA }o--|| PRESENTACION : "usa componente"
    
    PRECIO_HISTORICO {
        uint64 id PK
        uint64 presentacion_id FK
        decimal precio_anterior "precision(16,4)"
        decimal precio_nuevo "precision(16,4)"
        timestamp fecha_cambio
        uint64 usuario_id FK
        string motivo_cambio
        bool activo "default:true"
    }

    MONITOR_COCINA{
    uint64 id PK
    string monitor
    uint64 empresa_id FK "Empresa a la que pertenece"   
    }
    
    MONITOR_AUXILIAR{
    uint64 id PK
    string monitorauxiliar
    uint64 empresa_id FK "Empresa a la que pertenece"   
    }
    
    PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"   
        uint64 usuario_id
        uint64 tipo_producto_id FK       
        uint64 categoria_id FK
        uint64 area_despacho_id FK
        uint64 unidad_control_id FK
        uint64 impuesto_id FK 
        int64 monitor_cocina_id FK
        int64 monitor_auxiliar_id FK
        uint64 codigo_sunat_id    
        string nombre
        string descripcion
        string codigo UK
        string codigo_barras "Código de barras del producto"
        int tiempo_preparacion 
        bool es_plato "Se produce desde receta para vender?"
        bool es_mercaderia "Se compra para vender directamente?"
        bool es_insumo "Se compra para usar en recetas/transformación?"  
        bool es_semi_elaborado "Se produce desde receta para usar en otras recetas?"
        bool es_transformable "Este insumo puede ser procesado internamente en otros insumos/semi-elaborados?"
        bool es_compuesto 
        bool es_modificador 
        bool es_oferta 
        bool es_mensualidad
        bool tiene_controlstock
        bool tiene_controlvencimiento
        bool tiene_controllote
        bool es_produccion_diaria
        bool requiere_receta
        bool es_inactivo "No usable/visible globalmente"
    }

    PRESENTACION {
        uint64 id PK
        string nombre
        string descripcion  
        decimal cantidad_unidad_control "precision(16,4)"
        decimal precio "precision(16,4)"
        uint64 producto_id FK  
        bool es_inactivo "No usable/visible globalmente"
        bool es_presentacion_base       
        string codigo_barras "Código de barras de la presentación"
    }

    PRESENTACION_COMPUESTA {
        uint64 id PK
        uint64 presentacion_principal_id FK "Presentación del producto compuesto"
        uint64 presentacion_componente_id FK "Presentación que forma parte del compuesto"
        decimal cantidad "precision(16,4)"
        bool es_inactivo
        bool es_opcional
        int orden_visualizacion
    }
    PRESENTACION_PESO {
        uint64 id PK
        uint64 presentacion_id FK
        decimal precio_x_kilo "precision(16,4)"
        decimal peso_minimo "precision(16,4)"
        decimal peso_maximo "precision(16,4)"
        bool activo
    }

    PRESENTACION_TIEMPO {
        uint64 id PK
        uint64 presentacion_id FK
        decimal precio_x_hora "precision(16,4)"
        decimal tiempo_minimo "precision(8,2)"
        decimal tiempo_maximo "precision(8,2)"
        decimal fraccion_cobro "precision(8,2)"
        bool activo
    }

    PRODUCTO_MEDIA {
        uint64 id PK
        uint64 producto_id FK
        string url
        string tipo "IMAGEN|VIDEO"
        int orden

        bool es_principal
    }

    PRESENTACION_MEDIA {
        uint64 id PK
        uint64 presentacion_id FK
        string url
        string tipo "IMAGEN|VIDEO"
        int orden

        bool es_principal
    }

    TIPO_PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre

        string descripcion "Ayuda al usuario a entender el tipo"
    }

    AREA_DESPACHO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre        
        bool  imprimirticket
        string impresora_asociada         
    }

    CATEGORIA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre 
    }

    UNIDAD {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string simbolo      
        bool es_decimal
    }

    COMENTARIO_PRODUCTO {
        uint64 id PK
        string comentario
        uint64 producto_id FK "Null si es genérico"
        bool  esComentarioGenerico
    }

    GRUPO_CONTEO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        uint64 almacen_id FK "Referencia al ID del almacén en el microservicio de inventario"
        int orden_conteo "Orden para la secuencia de conteo"
    }

    CARTA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre 
        string descripcion
        string imagen
        int orden_visualizacion
        bool controlHorarioActivo
        time hora_inicio_servicio
        time hora_fin_servicio
    }
    
    FAMILIA {
        uint64 id PK
        uint64 carta_id FK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre        
        int orden_visualizacion
        string imagen
        bool activo "Indica si la familia está activa"
        bool destacado "Indica si la familia debe destacarse visualmente"
    }
    
    PRODUCTO_FAMILIA {
        uint64 id PK
        uint64 producto_id FK
        uint64 familia_id FK      
        int orden_visualizacion
        bool destacado       
        bool estadoDisponible
        bool controlHorarioActivo
        time hora_inicio_disponibilidad
        time hora_fin_disponibilidad
    }

    IMPUESTO {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string codigoSunat
        string nombre
        string descripcion
        decimal porcentaje "precision(16,4)"
    
    }

    PRODUCTO_GRUPO_CONTEO {
        uint64 id PK
        uint64 producto_id FK
        uint64 grupo_conteo_id FK
    }

    OFERTA {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        uint64 presentacion_id FK
        uint64 establecimiento_id FK
        string tipo_categoria
        string nombre
        string descripcion
        string imagen_url
        decimal precio_base "precision(16,4)"
        datetime fecha_inicio
        datetime fecha_fin
        bool tiene_control_fechas
        bool activo
        bool deleted
    }
    
    GRUPO_OFERTA {
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
    
    OFERTA_PRESENTACION {
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

    PRESENTACION_MODIFICADOR {
        uint64 id PK
        uint64 presentacion_id FK
        uint64 modificador_id FK
        bool required
        int orden_visualizacion
        bool activo
        bool deleted
    }

    MODIFICADOR {
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

    MODIFICADOR_OPCION {
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

    GRUPO_CONTEO ||--o{ PRODUCTO_GRUPO_CONTEO : "contiene"
:::