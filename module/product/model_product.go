package product

import (
	"time"

	"gorm.io/gorm"
)

// ListaPrecios representa la tabla LISTA_PRECIOS
type ListaPrecios struct {
	gorm.Model
	EmpresaID         uint64                      `json:"empresa_id" gorm:"index"`
	Nombre            string                      `json:"nombre" gorm:"type:varchar(100)"`
	Descripcion       string                      `json:"descripcion" gorm:"type:text"`
	EsPrecioPorMayor  bool                        `json:"es_precio_por_mayor" gorm:"default:false"`
	EsPrecioEspecial  bool                        `json:"es_precio_especial" gorm:"default:false"`
	ListaPresentacion []*ListaPreciosPresentacion `json:"lista_presentacion,omitempty" gorm:"foreignKey:ListaPreciosID"`
}

// ListaPreciosPresentacion representa la tabla LISTA_PRECIOS_PRESENTACION
type ListaPreciosPresentacion struct {
	gorm.Model
	ListaPreciosID uint64        `json:"lista_precios_id" gorm:"index"`
	PresentacionID uint64        `json:"presentacion_id" gorm:"index"`
	Precio         float64       `json:"precio" gorm:"type:decimal(16,4)"`
	ListaPrecios   *ListaPrecios `json:"lista_precios,omitempty" gorm:"foreignKey:ListaPreciosID"`
	Presentacion   *Presentacion `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// Categoria representa la tabla CATEGORIA
type Categoria struct {
	gorm.Model
	EmpresaID     uint64          `json:"empresa_id" gorm:"index"`
	Nombre        string          `json:"nombre" gorm:"type:varchar(100)"`
	SubCategorias []*SubCategoria `json:"sub_categorias,omitempty" gorm:"foreignKey:CategoriaID"`
}

// SubCategoria representa la tabla SUB_CATEGORIA
type SubCategoria struct {
	gorm.Model
	CategoriaID uint64      `json:"categoria_id" gorm:"index"`
	Nombre      string      `json:"nombre" gorm:"type:varchar(100)"`
	Activo      bool        `json:"activo" gorm:"default:true"`
	Categoria   *Categoria  `json:"categoria,omitempty" gorm:"foreignKey:CategoriaID"`
	Productos   []*Producto `json:"productos,omitempty" gorm:"foreignKey:SubCategoriaID"`
}

// Clase representa la tabla CLASE
type Clase struct {
	gorm.Model
	EmpresaID          uint64      `json:"empresa_id" gorm:"index"`
	Nombre             string      `json:"nombre" gorm:"type:varchar(100)"`
	Activo             bool        `json:"activo" gorm:"default:true"`
	OrdenVisualizacion int         `json:"orden_visualizacion"`
	Imagen                string  `json:"imagen" gorm:"type:varchar(255)"`
	Productos          []*Producto `json:"productos,omitempty" gorm:"foreignKey:ClaseID"`
}

type Producto struct {
	gorm.Model
	EmpresaID         uint64  `json:"empresa_id" gorm:"index"`
	UsuarioID         uint64  `json:"usuario_id" gorm:"index"`
	TipoProductoID    uint64  `json:"tipo_producto_id" gorm:"index"`           // Made optional
	SubCategoriaID    *uint64 `json:"sub_categoria_id,omitempty" gorm:"index"` // Made optional
	AreaDespachoID    *uint64 `json:"area_despacho_id,omitempty" gorm:"index"` // Made optional
	ClaseID           *uint64 `json:"clase_id,omitempty" gorm:"index"`         // Made optional as per ER diagram
	UnidadControlID   uint64  `json:"unidad_control_id" gorm:"index"`
	ImpuestoID        uint64  `json:"impuesto_id" gorm:"index"`
	MonitorCocinaID   *uint64 `json:"monitor_cocina_id,omitempty" gorm:"index"`   // Made optional
	MonitorAuxiliarID *uint64 `json:"monitor_auxiliar_id,omitempty" gorm:"index"` // Made optional
	CodigoSunatID     uint64  `json:"codigo_sunat_id" gorm:"index"`
	Nombre            string  `json:"nombre" gorm:"type:varchar(200)"`
	Descripcion       string  `json:"descripcion" gorm:"type:text"`
	Precio            float64 `json:"precio" gorm:"type:decimal(16,4)"`
	Codigo            string  `json:"codigo" gorm:"type:varchar(100);uniqueIndex"`
	CodigoBarras      string  `json:"codigo_barras" gorm:"type:varchar(100)"`
	TiempoPreparacion int     `json:"tiempo_preparacion"`

	EsPlato         bool `json:"es_plato" gorm:"default:false"`
	EsMercaderia    bool `json:"es_mercaderia" gorm:"default:false"`
	EsInsumo        bool `json:"es_insumo" gorm:"default:false"`
	EsSemiElaborado bool `json:"es_semi_elaborado" gorm:"default:false"`
	EsCompuesto     bool `json:"es_compuesto" gorm:"default:false"`
	EsOferta        bool `json:"es_oferta" gorm:"default:false"`
	EsMensualidad   bool `json:"es_mensualidad" gorm:"default:false"`
	EsSuscripcion   bool `json:"es_suscripcion" gorm:"default:false"`

	TpvVisible                bool `json:"tpv_visible" gorm:"default:true"`
	EsTransformable           bool `json:"es_transformable" gorm:"default:false"`
	EsModificador             bool `json:"es_modificador" gorm:"default:false"`
	TieneControlStock         bool `json:"tiene_controlstock" gorm:"default:false"`
	TieneControlVencimiento   bool `json:"tiene_controlvencimiento" gorm:"default:false"`
	TieneControlLote          bool `json:"tiene_controllote" gorm:"default:false"`
	TienePeso                 bool `json:"tiene_peso" gorm:"default:false"`
	TieneTiempo               bool `json:"tiene_tiempo" gorm:"default:false"`
	TieneDescuentoPorCantidad bool `json:"tiene_descuento_por_cantidad" gorm:"default:false"`
	EsProduccionDiaria        bool `json:"es_produccion_diaria" gorm:"default:false"`
	RequiereReceta            bool `json:"requiere_receta" gorm:"default:false"`
	EsInactivo                bool `json:"es_inactivo" gorm:"default:false"`
	AlmacenVisible            bool `json:"almacen_visible" gorm:"default:false"`

	// Relaciones
	TipoProducto          *TipoProducto                 `json:"tipo_producto,omitempty" gorm:"foreignKey:TipoProductoID"`
	SubCategoria          *SubCategoria                 `json:"sub_categoria,omitempty" gorm:"foreignKey:SubCategoriaID"` // Changed from Categoria
	AreaDespacho          *AreaDespacho                 `json:"area_despacho,omitempty" gorm:"foreignKey:AreaDespachoID"`
	UnidadControl         *Unidad                       `json:"unidad_control,omitempty" gorm:"foreignKey:UnidadControlID"`
	Impuesto              *Impuesto                     `json:"impuesto,omitempty" gorm:"foreignKey:ImpuestoID"`
	MonitorCocina         *MonitorCocina                `json:"monitor_cocina,omitempty" gorm:"foreignKey:MonitorCocinaID"`
	MonitorAuxiliar       *MonitorAuxiliar              `json:"monitor_auxiliar,omitempty" gorm:"foreignKey:MonitorAuxiliarID"`
	Presentaciones        []*Presentacion               `json:"presentaciones,omitempty" gorm:"foreignKey:ProductoID"`
	ComentariosProducto   []*ComentarioProducto         `json:"comentarios_producto,omitempty" gorm:"foreignKey:ProductoID"`
	ProductoGrupoConteo   []*ProductoGrupoConteo        `json:"producto_grupo_conteo,omitempty" gorm:"foreignKey:ProductoID"`
	ProductoMedia         []*ProductoMedia              `json:"producto_media,omitempty" gorm:"foreignKey:ProductoID"`
	PreciosDescuento      []*PrecioDescuentoPorCantidad `json:"precios_descuento,omitempty" gorm:"foreignKey:ProductoID"`
	ProductoFamilias      []*ProductoFamilia            `json:"producto_familias,omitempty" gorm:"foreignKey:ProductoID"`
	ProductoModificadores []*ProductoModificador        `json:"producto_modificadores,omitempty" gorm:"foreignKey:ProductoID"`
	Clase                 *Clase                        `json:"clase,omitempty" gorm:"foreignKey:ClaseID"`
}

// Presentacion representa la tabla PRESENTACION
type Presentacion struct {
	gorm.Model
	ProductoID            uint64  `json:"producto_id" gorm:"index"`
	Nombre                string  `json:"nombre" gorm:"type:varchar(200)"`
	Descripcion           string  `json:"descripcion" gorm:"type:text"`
	CantidadUnidadControl float64 `json:"cantidad_unidad_control" gorm:"type:decimal(16,4)"`
	Precio                float64 `json:"precio" gorm:"type:decimal(16,4)"`
	EsInactivo            bool    `json:"es_inactivo" gorm:"default:false"`
	EsPresentacionBase    bool    `json:"es_presentacion_base" gorm:"default:false"`
	CodigoBarras          string  `json:"codigo_barras" gorm:"type:varchar(100)"`
	Imagen                string  `json:"imagen" gorm:"type:varchar(255)"`
	EsCompra              bool    `json:"es_compra" gorm:"default:false"`

	// Relaciones
	Producto                 *Producto                `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
	PresentacionPeso         []*PresentacionPeso      `json:"presentacion_peso,omitempty" gorm:"foreignKey:PresentacionID"`
	PresentacionTiempo       []*PresentacionTiempo    `json:"presentacion_tiempo,omitempty" gorm:"foreignKey:PresentacionID"`
	PreciosHistoricos        []*PrecioHistorico       `json:"precios_historicos,omitempty" gorm:"foreignKey:PresentacionID"`
	PresentacionesCompuestas []*PresentacionCompuesta `json:"presentaciones_compuestas,omitempty" gorm:"foreignKey:PresentacionPrincipalID"`
	ComponentesCompuestos    []*PresentacionCompuesta `json:"componentes_compuestos,omitempty" gorm:"foreignKey:PresentacionComponenteID"`
	OfertaPresentaciones     []*OfertaPresentacion    `json:"oferta_presentaciones,omitempty" gorm:"foreignKey:PresentacionID"`
}

// TipoProducto representa la tabla TIPO_PRODUCTO
type TipoProducto struct {
	gorm.Model
	EmpresaID   uint64 `json:"empresa_id" gorm:"index"`
	Nombre      string `json:"nombre" gorm:"type:varchar(100)"`
	Descripcion string `json:"descripcion" gorm:"type:text"`

	EsPlato         bool `json:"es_plato" gorm:"default:false"`
	EsMercaderia    bool `json:"es_mercaderia" gorm:"default:false"`
	EsInsumo        bool `json:"es_insumo" gorm:"default:false"`
	EsSemiElaborado bool `json:"es_semi_elaborado" gorm:"default:false"`
	EsCompuesto     bool `json:"es_compuesto" gorm:"default:false"`
	EsOferta        bool `json:"es_oferta" gorm:"default:false"`
	EsMensualidad   bool `json:"es_mensualidad" gorm:"default:false"`
	EsSuscripcion   bool `json:"es_suscripcion" gorm:"default:false"`

	TpvVisible                bool        `json:"tpv_visible" gorm:"default:true"`
	AlmacenVisible            bool        `json:"almacen_visible" gorm:"default:false"`
	EsTransformable           bool        `json:"es_transformable" gorm:"default:false"`
	EsModificador             bool        `json:"es_modificador" gorm:"default:false"`
	TieneControlStock         bool        `json:"tiene_controlstock" gorm:"default:false"`
	TieneControlVencimiento   bool        `json:"tiene_controlvencimiento" gorm:"default:false"`
	TieneControlLote          bool        `json:"tiene_controllote" gorm:"default:false"`
	TienePeso                 bool        `json:"tiene_peso" gorm:"default:false"`
	TieneTiempo               bool        `json:"tiene_tiempo" gorm:"default:false"`
	TieneDescuentoPorCantidad bool        `json:"tiene_descuento_por_cantidad" gorm:"default:false"`
	EsProduccionDiaria        bool        `json:"es_produccion_diaria" gorm:"default:false"`
	RequiereReceta            bool        `json:"requiere_receta" gorm:"default:false"`
	EsInactivo                bool        `json:"es_inactivo" gorm:"default:false"`
	EsTipoUniversal           bool        `json:"es_tipo_universal" gorm:"default:false"`
	RequiereTipoProducto      bool        `json:"requiere_tipo_producto"`
	RequiereSubcategoria      bool        `json:"requiere_subcategoria"`
	RequiereAreaDespacho      bool        `json:"requiere_area_despacho"`
	RequiereMonitorCocina     bool        `json:"requiere_monitor_cocina"`
	RequiereMonitorAuxiliar   bool        `json:"requiere_monitor_auxiliar"`
	RequiereClase             bool        `json:"requiere_clase"`
	Imagen                    string      `json:"imagen" gorm:"type:varchar(255)"`
	Productos                 []*Producto `json:"productos,omitempty" gorm:"foreignKey:TipoProductoID"`
}

// AreaDespacho representa la tabla AREA_DESPACHO
type AreaDespacho struct {
	gorm.Model
	EmpresaID         uint64      `json:"empresa_id" gorm:"index"`
	Nombre            string      `json:"nombre" gorm:"type:varchar(100)"`
	ImprimirTicket    bool        `json:"imprimir_ticket" gorm:"default:false"`
	ImpresoraAsociada string      `json:"impresora_asociada" gorm:"type:varchar(100)"`
	Productos         []*Producto `json:"productos,omitempty" gorm:"foreignKey:AreaDespachoID"`
}

// Unidad representa la tabla UNIDAD
type Unidad struct {
	gorm.Model
	EmpresaID uint64      `json:"empresa_id" gorm:"index"`
	Nombre    string      `json:"nombre" gorm:"type:varchar(100)"`
	Simbolo   string      `json:"simbolo" gorm:"type:varchar(10)"`
	EsDecimal bool        `json:"es_decimal" gorm:"default:false"`
	Productos []*Producto `json:"productos,omitempty" gorm:"foreignKey:UnidadControlID"`
}

// ComentarioProducto representa la tabla COMENTARIO_PRODUCTO
type ComentarioProducto struct {
	gorm.Model
	ProductoID           *uint64   `json:"producto_id,omitempty" gorm:"index"` // Changed to pointer type *uint64
	Comentario           string    `json:"comentario" gorm:"type:text"`
	EsComentarioGenerico bool      `json:"es_comentario_generico" gorm:"default:false"`
	Producto             *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
}

// GrupoConteo representa la tabla GRUPO_CONTEO
type GrupoConteo struct {
	gorm.Model
	EmpresaID      uint64                 `json:"empresa_id" gorm:"index"`
	Nombre         string                 `json:"nombre" gorm:"type:varchar(100)"`
	Descripcion    string                 `json:"descripcion" gorm:"type:text"`
	AlmacenID      uint64                 `json:"almacen_id" gorm:"index"`
	OrdenConteo    int                    `json:"orden_conteo"`
	ProductoGrupos []*ProductoGrupoConteo `json:"producto_grupos,omitempty" gorm:"foreignKey:GrupoConteoID"`
}

// ProductoGrupoConteo representa la tabla PRODUCTO_GRUPO_CONTEO
type ProductoGrupoConteo struct {
	gorm.Model
	ProductoID    uint64       `json:"producto_id" gorm:"index"`
	GrupoConteoID uint64       `json:"grupo_conteo_id" gorm:"index"`
	Producto      *Producto    `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
	GrupoConteo   *GrupoConteo `json:"grupo_conteo,omitempty" gorm:"foreignKey:GrupoConteoID"`
}

// Carta representa la tabla CARTA
type Carta struct {
	gorm.Model
	EmpresaID            uint64     `json:"empresa_id" gorm:"index"`
	Nombre               string     `json:"nombre" gorm:"type:varchar(100)"`
	Descripcion          string     `json:"descripcion" gorm:"type:text"`
	Imagen               string     `json:"imagen" gorm:"type:varchar(255)"`
	OrdenVisualizacion   int        `json:"orden_visualizacion"`
	ControlHorarioActivo bool       `json:"control_horario_activo" gorm:"default:false"`
	HoraInicioServicio   time.Time  `json:"hora_inicio_servicio"`
	HoraFinServicio      time.Time  `json:"hora_fin_servicio"`
	Familias             []*Familia `json:"familias,omitempty" gorm:"foreignKey:CartaID"`
}

// Familia representa la tabla FAMILIA
type Familia struct {
	gorm.Model
	CartaID            uint64             `json:"carta_id" gorm:"index"`
	EmpresaID          uint64             `json:"empresa_id" gorm:"index"`
	Nombre             string             `json:"nombre" gorm:"type:varchar(100)"`
	OrdenVisualizacion int                `json:"orden_visualizacion"`
	Imagen             string             `json:"imagen" gorm:"type:varchar(255)"`
	Activo             bool               `json:"activo" gorm:"default:true"`
	Destacado          bool               `json:"destacado" gorm:"default:false"`
	Carta              *Carta             `json:"carta,omitempty" gorm:"foreignKey:CartaID"`
	ProductoFamilias   []*ProductoFamilia `json:"producto_familias,omitempty" gorm:"foreignKey:FamiliaID"`
}

// ProductoFamilia representa la tabla PRODUCTO_FAMILIA
type ProductoFamilia struct {
	gorm.Model
	ProductoID               uint64    `json:"producto_id" gorm:"index"`
	FamiliaID                uint64    `json:"familia_id" gorm:"index"`
	OrdenVisualizacion       int       `json:"orden_visualizacion"`
	Destacado                bool      `json:"destacado" gorm:"default:false"`
	EstadoDisponible         bool      `json:"estado_disponible" gorm:"default:true"`
	ControlHorarioActivo     bool      `json:"control_horario_activo" gorm:"default:false"`
	HoraInicioDisponibilidad time.Time `json:"hora_inicio_disponibilidad"`
	HoraFinDisponibilidad    time.Time `json:"hora_fin_disponibilidad"`
	Producto                 *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
	Familia                  *Familia  `json:"familia,omitempty" gorm:"foreignKey:FamiliaID"`
}

// ProductoMedia representa la tabla PRODUCTO_MEDIA
type ProductoMedia struct {
	gorm.Model
	ProductoID  uint64    `json:"producto_id" gorm:"index"`
	URL         string    `json:"url" gorm:"type:varchar(255)"`
	Tipo        string    `json:"tipo" gorm:"type:varchar(10)"` // IMAGEN|VIDEO
	Orden       int       `json:"orden"`
	EsPrincipal bool      `json:"es_principal" gorm:"default:false"`
	Producto    *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
}

// PrecioHistorico representa la tabla PRECIO_HISTORICO
type PrecioHistorico struct {
	gorm.Model
	PresentacionID uint64        `json:"presentacion_id" gorm:"index"`
	PrecioAnterior float64       `json:"precio_anterior" gorm:"type:decimal(16,4)"`
	PrecioNuevo    float64       `json:"precio_nuevo" gorm:"type:decimal(16,4)"`
	UsuarioID      uint64        `json:"usuario_id" gorm:"index"`
	PrecioActivo   bool          `json:"precio_activo" gorm:"default:true"`
	Presentacion   *Presentacion `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// PrecioDescuentoPorCantidad representa la tabla PRECIO_DESCUENTO_POR_CANTIDAD
type PrecioDescuentoPorCantidad struct {
	gorm.Model
	ProductoID        uint64    `json:"producto_id" gorm:"index"`
	CantidadMinima    int       `json:"cantidad_minima"`
	CantidadMaxima    int       `json:"cantidad_maxima"`
	UsaPrecio         bool      `json:"usa_precio" gorm:"default:true"`
	PrecioAlternativo float64   `json:"precio_alternativo" gorm:"type:decimal(16,4)"`
	Descuento         float64   `json:"descuento" gorm:"type:decimal(16,4)"`
	Producto          *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
}

// MonitorCocina representa la tabla MONITOR_COCINA
type MonitorCocina struct {
	gorm.Model
	EmpresaID uint64      `json:"empresa_id" gorm:"index"`
	Monitor   string      `json:"monitor" gorm:"type:varchar(100)"`
	ImprimirTicket    bool        `json:"imprimir_ticket" gorm:"default:false"`
	ImpresoraAsociada string      `json:"impresora_asociada" gorm:"type:varchar(100)"`
	Productos []*Producto `json:"productos,omitempty" gorm:"foreignKey:MonitorCocinaID"`
}

// MonitorAuxiliar representa la tabla MONITOR_AUXILIAR
type MonitorAuxiliar struct {
	gorm.Model
	EmpresaID       uint64      `json:"empresa_id" gorm:"index"`
	MonitorAuxiliar string      `json:"monitor_auxiliar" gorm:"type:varchar(100)"`
	Productos       []*Producto `json:"productos,omitempty" gorm:"foreignKey:MonitorAuxiliarID"`
}

// Impuesto representa la tabla IMPUESTO
type Impuesto struct {
	gorm.Model
	EmpresaID      uint64      `json:"empresa_id" gorm:"index"`
	NombreImpuesto string      `json:"nombre_impuesto" gorm:"type:varchar(100)"`
	Porcentaje     float64     `json:"porcentaje" gorm:"type:decimal(16,4)"`
	IdeTributo     uint64      `json:"ide_tributo" gorm:"index"`
	NomTributo     string      `json:"nom_tributo" gorm:"type:varchar(100)"`
	CodTipTributo  string      `json:"cod_tip_tributo" gorm:"type:varchar(20)"`
	Productos      []*Producto `json:"productos,omitempty" gorm:"foreignKey:ImpuestoID"`
}

// Oferta representa la tabla OFERTA
type Oferta struct {
	gorm.Model
	EmpresaID               uint64         `json:"empresa_id" gorm:"index"`
	PresentacionID          uint64         `json:"presentacion_id" gorm:"index"`
	FechaInicio             time.Time      `json:"fecha_inicio"`
	FechaFin                time.Time      `json:"fecha_fin"`
	StockOfertado           float64        `json:"stock_ofertado" gorm:"type:decimal(16,4)"`
	StockActual             float64        `json:"stock_actual" gorm:"type:decimal(16,4)"`
	ControlCantidadOfertada bool           `json:"control_cantidad_ofertada" gorm:"default:false"`
	TieneControlFechas      bool           `json:"tiene_control_fechas" gorm:"default:false"`
	Activo                  bool           `json:"activo" gorm:"default:true"`
	Deleted                 bool           `json:"deleted" gorm:"default:false"`
	GruposOferta            []*GrupoOferta `json:"grupos_oferta,omitempty" gorm:"foreignKey:OfertaID"`
	Presentacion            *Presentacion  `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// GrupoOferta representa la tabla GRUPO_OFERTA
type GrupoOferta struct {
	gorm.Model
	OfertaID             uint64                `json:"oferta_id" gorm:"index"`
	Nombre               string                `json:"nombre" gorm:"type:varchar(100)"`
	MinSelecciones       uint8                 `json:"min_selecciones"`
	MaxSelecciones       uint8                 `json:"max_selecciones"`
	PrecioBase           float64               `json:"precio_base" gorm:"type:decimal(16,4)"`
	OrdenVisualizacion   int                   `json:"orden_visualizacion"`
	Activo               bool                  `json:"activo" gorm:"default:true"`
	Deleted              bool                  `json:"deleted" gorm:"default:false"`
	Oferta               *Oferta               `json:"oferta,omitempty" gorm:"foreignKey:OfertaID"`
	OfertaPresentaciones []*OfertaPresentacion `json:"oferta_presentaciones,omitempty" gorm:"foreignKey:GrupoOfertaID"`
}

// OfertaPresentacion representa la tabla OFERTA_PRESENTACION
type OfertaPresentacion struct {
	gorm.Model
	OfertaID        uint64        `json:"oferta_id" gorm:"index"`
	GrupoOfertaID   uint64        `json:"grupo_oferta_id" gorm:"index"`
	PresentacionID  uint64        `json:"presentacion_id" gorm:"index"`
	PrecioAdicional float64       `json:"precio_adicional" gorm:"type:decimal(16,4)"`
	Activo          bool          `json:"activo" gorm:"default:true"`
	Deleted         bool          `json:"deleted" gorm:"default:false"`
	Oferta          *Oferta       `json:"oferta,omitempty" gorm:"foreignKey:OfertaID"`
	GrupoOferta     *GrupoOferta  `json:"grupo_oferta,omitempty" gorm:"foreignKey:GrupoOfertaID"`
	Presentacion    *Presentacion `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// PresentacionCompuesta representa la tabla PRESENTACION_COMPUESTA
type PresentacionCompuesta struct {
	gorm.Model
	PresentacionPrincipalID  uint64        `json:"presentacion_principal_id" gorm:"index"`
	PresentacionComponenteID uint64        `json:"presentacion_componente_id" gorm:"index"`
	Cantidad                 float64       `json:"cantidad" gorm:"type:decimal(16,4)"`
	EsInactivo               bool          `json:"es_inactivo" gorm:"default:false"`
	EsOpcional               bool          `json:"es_opcional" gorm:"default:false"`
	OrdenVisualizacion       int           `json:"orden_visualizacion"`
	PresentacionPrincipal    *Presentacion `json:"presentacion_principal,omitempty" gorm:"foreignKey:PresentacionPrincipalID"`
	PresentacionComponente   *Presentacion `json:"presentacion_componente,omitempty" gorm:"foreignKey:PresentacionComponenteID"`
}

// PresentacionPeso representa la tabla PRESENTACION_PESO
type PresentacionPeso struct {
	gorm.Model
	PresentacionID uint64        `json:"presentacion_id" gorm:"index"`
	PrecioXKilo    float64       `json:"precio_x_kilo" gorm:"type:decimal(16,4)"`
	PesoMinimo     float64       `json:"peso_minimo" gorm:"type:decimal(16,4)"`
	PesoMaximo     float64       `json:"peso_maximo" gorm:"type:decimal(16,4)"`
	Activo         bool          `json:"activo" gorm:"default:true"`
	Presentacion   *Presentacion `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// PresentacionTiempo representa la tabla PRESENTACION_TIEMPO
type PresentacionTiempo struct {
	gorm.Model
	PresentacionID uint64        `json:"presentacion_id" gorm:"index"`
	PrecioXHora    float64       `json:"precio_x_hora" gorm:"type:decimal(16,4)"`
	TiempoMinimo   float64       `json:"tiempo_minimo" gorm:"type:decimal(8,2)"`
	TiempoMaximo   float64       `json:"tiempo_maximo" gorm:"type:decimal(8,2)"`
	FraccionCobro  float64       `json:"fraccion_cobro" gorm:"type:decimal(8,2)"`
	Activo         bool          `json:"activo" gorm:"default:true"`
	Presentacion   *Presentacion `json:"presentacion,omitempty" gorm:"foreignKey:PresentacionID"`
}

// Modificador representa la tabla MODIFICADOR
type Modificador struct {
	gorm.Model
	EmpresaID             uint64                 `json:"empresa_id" gorm:"index"`
	Nombre                string                 `json:"nombre" gorm:"type:varchar(100)"`
	MinSelecciones        uint8                  `json:"min_selecciones"`
	MaxSelecciones        uint8                  `json:"max_selecciones"`
	MultipleSeleccion     bool                   `json:"multiple_seleccion" gorm:"default:false"`
	Activo                bool                   `json:"activo" gorm:"default:true"`
	Deleted               bool                   `json:"deleted" gorm:"default:false"`
	Opciones              []*ModificadorOpcion   `json:"opciones,omitempty" gorm:"foreignKey:ModificadorID"`
	ProductoModificadores []*ProductoModificador `json:"producto_modificadores,omitempty" gorm:"foreignKey:ModificadorID"`
}

// ModificadorOpcion representa la tabla MODIFICADOR_OPCION
type ModificadorOpcion struct {
	gorm.Model
	ModificadorID      uint64       `json:"modificador_id" gorm:"index"`
	Nombre             string       `json:"nombre" gorm:"type:varchar(100)"`
	PrecioAdicional    float64      `json:"precio_adicional" gorm:"type:decimal(16,4)"`
	OrdenVisualizacion int          `json:"orden_visualizacion"`
	Activo             bool         `json:"activo" gorm:"default:true"`
	Deleted            bool         `json:"deleted" gorm:"default:false"`
	ImagenURL          string       `json:"imagen_url" gorm:"type:varchar(255)"`
	Modificador        *Modificador `json:"modificador,omitempty" gorm:"foreignKey:ModificadorID"`
}

// ProductoModificador representa la tabla PRODUCTO_MODIFICADOR
type ProductoModificador struct {
	gorm.Model
	ProductoID         uint64       `json:"producto_id" gorm:"index"`
	ModificadorID      uint64       `json:"modificador_id" gorm:"index"`
	Required           bool         `json:"required" gorm:"default:false"`
	OrdenVisualizacion int          `json:"orden_visualizacion"`
	Activo             bool         `json:"activo" gorm:"default:true"`
	Deleted            bool         `json:"deleted" gorm:"default:false"`
	Producto           *Producto    `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
	Modificador        *Modificador `json:"modificador,omitempty" gorm:"foreignKey:ModificadorID"`
}
