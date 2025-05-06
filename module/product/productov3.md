::: mermaid
erDiagram
    %% Update relationships to match model_product.go
    EMPRESA ||--o{ PRODUCTO : tiene
    EMPRESA ||--o{ CATEGORIA : tiene   
    EMPRESA ||--o{ CLASE : tiene
    CATEGORIA ||--o{ SUB_CATEGORIA : contiene
    SUB_CATEGORIA ||--o{ PRODUCTO : tiene
    EMPRESA ||--o{ TIPO_PRODUCTO : tiene
    EMPRESA ||--o{ AREA_DESPACHO : tiene
    EMPRESA ||--o{ UNIDAD : tiene
    EMPRESA ||--o{ CARTA : tiene   
    EMPRESA ||--o{ GRUPO_CONTEO : tiene
    EMPRESA ||--o{ MONITOR_COCINA : tiene
    EMPRESA ||--o{ MONITOR_AUXILIAR : tiene
    EMPRESA ||--o{ OFERTA : tiene
    EMPRESA ||--o{ MODIFICADOR : tiene
    
    CLASE ||--o{ PRODUCTO : "clasifica" 
    
    PRODUCTO ||--o{ PRESENTACION : tiene   
    PRODUCTO }o--|| TIPO_PRODUCTO : clasifica
    PRODUCTO }o--|| AREA_DESPACHO : despacha
    PRODUCTO }o--|| UNIDAD : controla
    PRODUCTO ||--o{ COMENTARIO_PRODUCTO : permite
    PRODUCTO ||--o{ PRODUCTO_GRUPO_CONTEO : "puede tener"
    PRODUCTO ||--o{ PRODUCTO_MEDIA : "tiene"    
    CARTA ||--|{ FAMILIA : "contiene"
    FAMILIA ||--|{ PRODUCTO_FAMILIA : "contiene"
    PRODUCTO ||--|{ PRODUCTO_FAMILIA : "está en"
    PRODUCTO }o--|| IMPUESTO : "puede tener"
    PRODUCTO }o--|| MONITOR_COCINA : "usa"
    PRODUCTO }o--|| MONITOR_AUXILIAR : "usa"
    OFERTA ||--o{ GRUPO_OFERTA : "contiene"
    GRUPO_OFERTA ||--|{ OFERTA_PRESENTACION : "organiza"
    OFERTA_PRESENTACION }|--|| PRESENTACION : "usa" 
    PRODUCTO ||--o{ PRODUCTO_MODIFICADOR : "tiene"
    PRODUCTO_MODIFICADOR }o--|| MODIFICADOR : "usa"
    MODIFICADOR ||--|{ MODIFICADOR_OPCION : "tiene"
    PRESENTACION |o--|| OFERTA : "puede ser"
    PRESENTACION ||--o{ PRESENTACION_PESO : "puede tener"
    PRESENTACION ||--o{ PRESENTACION_TIEMPO : "puede tener"
    PRESENTACION ||--o{ PRECIO_HISTORICO : "registra cambios"
    PRESENTACION ||--o{ PRESENTACION_COMPUESTA : "contiene"
    PRESENTACION_COMPUESTA }o--|| PRESENTACION : "usa componente"
    PRODUCTO ||--o{ PRECIO_DESCUENTO_POR_CANTIDAD : "configura"
    
    LISTA_PRECIOS ||--o{ LISTA_PRECIOS_PRESENTACION : "asignada"

    LISTA_PRECIOS {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre
        string descripcion
        bool es_precio_por_mayor
        bool es_precio_especial
    }

    LISTA_PRECIOS_PRESENTACION {
        uint64 id PK
        uint64 lista_precios_id FK
        uint64 presentacion_id FK
        decimal precio "precision(16,4)"
    }

    PRECIO_HISTORICO {
        uint64 id PK
        uint64 presentacion_id FK
        decimal precio_anterior "precision(16,4)"
        decimal precio_nuevo "precision(16,4)"
        uint64 usuario_id FK
        bool precioactivo "default:true"
    }

    PRECIO_DESCUENTO_POR_CANTIDAD {
        uint64 id PK
        uint64 producto_id FK
        int cantidad_minima
        int cantidad_maxima
        bool usa_precio "default true"
        decimal precio_alternativo "precision(16,4)"
        decimal descuento "precision(16,4)"
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
        uint64 empresa_id FK
        uint64 usuario_id FK
        uint64 tipo_producto_id FK "optional"
        uint64 sub_categoria_id FK "optional"
        uint64 area_despacho_id FK "optional"
        uint64 unidad_control_id FK
        uint64 impuesto_id FK
        uint64 monitor_cocina_id FK "optional"
        uint64 monitor_auxiliar_id FK "optional"
        uint64 clase_id FK "optional"
        uint64 codigo_sunat_id FK
        string nombre "varchar(200)"
        string descripcion "text"
        decimal precio "precision(16,4)"
        string codigo "varchar(100),unique"
        string codigo_barras "varchar(100)"
        int tiempo_preparacion      

        bool es_plato "default:false"
        bool es_mercaderia "default:false"
        bool es_insumo "default:false"
        bool es_semi_elaborado "default:false"       
        bool es_compuesto "default:false"        
        bool es_oferta "default:false" 
        bool es_mensualidad "default:false"
	    bool es_suscripcion "default:false"

        bool tpv_visible "default:true"
        bool almacen_visible "default:false"
        bool es_transformable "default:false"
        bool es_modificador "default:false"
        bool tiene_controlstock "default:false"
        bool tiene_controlvencimiento "default:false"
        bool tiene_controllote "default:false"
        bool tiene_peso "default:false"
        bool tiene_tiempo "default:false"
        bool tiene_descuento_por_cantidad "default:false"
        bool es_produccion_diaria "default:false"
        bool requiere_receta "default:false"
        bool es_inactivo "default:false"
        
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
        string url
        bool es_compra "default:false"
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

    TIPO_PRODUCTO {
        uint64 id PK
        uint64 empresa_id FK "Opcional - Solo si es personalizado"
        string nombre
        string descripcion "Ayuda al usuario a entender el tipo"

        bool es_plato "default:false"
        bool es_mercaderia "default:false"
        bool es_insumo "default:false"
        bool es_semi_elaborado "default:false"       
        bool es_compuesto "default:false"        
        bool es_oferta "default:false" 
        bool es_mensualidad "default:false"
	    bool es_suscripcion "default:false"

        bool tpv_visible "default:true"
        bool almacen_visible "default:false"
        bool es_transformable "default:false"
        bool es_modificador "default:false"
        bool tiene_controlstock "default:false"
        bool tiene_controlvencimiento "default:false"
        bool tiene_controllote "default:false"
        bool tiene_peso "default:false"
        bool tiene_tiempo "default:false"
        bool tiene_descuento_por_cantidad "default:false"
        bool es_produccion_diaria "default:false"
        bool requiere_receta "default:false"
        bool es_inactivo "default:false"
        bool es_tipo_universal "default:false"
        bool requiere_tipo_producto "Hace visible tipo_producto_id"
        bool requiere_subcategoria "Hace visible sub_categoria_id"
        bool requiere_area_despacho "Hace visible area_despacho_id"
        bool requiere_monitor_cocina "Hace visible monitor_cocina_id"
        bool requiere_monitor_auxiliar "Hace visible monitor_auxiliar_id"
        bool requiere_clase "Hace visible clase_id"
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

    SUB_CATEGORIA {
        uint64 id PK
        uint64 categoria_id FK
        string nombre "varchar(100)"
        bool activo "default:true"
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
        string nombreimpuesto       
        decimal porcentaje "precision(16,4)"
        uint64 ideTributo
        string nomTributo
        string codTipTributo    
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
        datetime fecha_inicio
        datetime fecha_fin
        decimal stock_ofertado "precision(16,4)"
        decimal stock_actual "precision(16,4)"
        bool control_cantidad_ofertada "por defecto false"
        bool tiene_control_fechas "por defecto false"
        bool activo
        bool deleted
    }
    
    GRUPO_OFERTA {
        uint64 id PK
        uint64 oferta_id FK
        string nombre
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
        bool activo
        bool deleted
    }
    PRODUCTO_MODIFICADOR {
        uint64 id PK
        uint64 producto_id FK
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
        decimal precio_adicional "precision(16,4)"
        int orden_visualizacion
        bool activo
        bool deleted
        string imagen_url 
    }

    CLASE {
        uint64 id PK
        uint64 empresa_id FK "Empresa a la que pertenece"
        string nombre "varchar(100)"
        bool activo "default:true"
        int orden_visualizacion
    }

    GRUPO_CONTEO ||--o{ PRODUCTO_GRUPO_CONTEO : "contiene"
:::