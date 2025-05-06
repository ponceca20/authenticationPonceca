package product

import (
	"fmt"
	"time"

	"practicev2/database"
	"practicev2/registry"

	"github.com/gofiber/fiber/v2"
)

// Seeders siembra datos para todas las tablas definidas en model_product.go siguiendo el orden correcto de relaciones.
func Seeders() bool {
	// ...existing code: conectar a base de datos...
	database.ConnectDatabase()
	db := database.DBconn

	// 1. Sembrar tablas independientes
	tiposProducto := []TipoProducto{
		{EmpresaID: 1, Nombre: "PLATO", Descripcion: "Tipo de producto PLATO",
			EsPlato:    true,
			TpvVisible: true, AlmacenVisible: false, EsTransformable: false, EsModificador: false,
			TieneControlStock: false, TieneControlVencimiento: false, TieneControlLote: false,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: false, RequiereReceta: true, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: false, RequiereSubcategoria: true, RequiereAreaDespacho: true,
			RequiereMonitorCocina: true, RequiereMonitorAuxiliar: false, RequiereClase: true,
		},
		{EmpresaID: 1, Nombre: "MERCADERIA", Descripcion: "Tipo de producto MERCADERIA",
			EsMercaderia: true,
			TpvVisible:   true, AlmacenVisible: true, EsTransformable: false, EsModificador: false,
			TieneControlStock: true, TieneControlVencimiento: false, TieneControlLote: true,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: false, RequiereReceta: false, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: false, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "INSUMO", Descripcion: "Tipo de producto INSUMO",
			EsInsumo:   true,
			TpvVisible: false, AlmacenVisible: true, EsTransformable: false, EsModificador: false,
			TieneControlStock: true, TieneControlVencimiento: true, TieneControlLote: true,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: true, RequiereReceta: false, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: false, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "SEMI_ELABORADO", Descripcion: "Tipo de producto SEMI_ELABORADO",
			EsSemiElaborado: true,
			TpvVisible:      false, AlmacenVisible: true, EsTransformable: true, EsModificador: false,
			TieneControlStock: true, TieneControlVencimiento: true, TieneControlLote: true,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: true, RequiereReceta: true, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: true, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "COMPUESTO", Descripcion: "Tipo de producto COMPUESTO",
			EsCompuesto: true,
			TpvVisible:  true, AlmacenVisible: false, EsTransformable: true, EsModificador: false,
			TieneControlStock: false, TieneControlVencimiento: false, TieneControlLote: false,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: false, RequiereReceta: true, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: true, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "OFERTA", Descripcion: "Tipo de producto OFERTA",
			EsOferta:   true,
			TpvVisible: true, AlmacenVisible: false, EsTransformable: false, EsModificador: false,
			TieneControlStock: false, TieneControlVencimiento: false, TieneControlLote: false,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: true,
			EsProduccionDiaria: false, RequiereReceta: false, EsInactivo: false, EsTipoUniversal: false,
			RequiereTipoProducto: false, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "MENSUALIDAD", Descripcion: "Tipo de producto MENSUALIDAD",
			EsMensualidad: true,
			TpvVisible:    false, AlmacenVisible: false, EsTransformable: false, EsModificador: false,
			TieneControlStock: false, TieneControlVencimiento: false, TieneControlLote: false,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: false, RequiereReceta: false, EsInactivo: false, EsTipoUniversal: true,
			RequiereTipoProducto: false, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
		{EmpresaID: 1, Nombre: "SUSCRIPCION", Descripcion: "Tipo de producto SUSCRIPCION",
			EsSuscripcion: true,
			TpvVisible:    false, AlmacenVisible: false, EsTransformable: false, EsModificador: false,
			TieneControlStock: false, TieneControlVencimiento: false, TieneControlLote: false,
			TienePeso: false, TieneTiempo: false, TieneDescuentoPorCantidad: false,
			EsProduccionDiaria: false, RequiereReceta: false, EsInactivo: false, EsTipoUniversal: true,
			RequiereTipoProducto: false, RequiereSubcategoria: false, RequiereAreaDespacho: false,
			RequiereMonitorCocina: false, RequiereMonitorAuxiliar: false, RequiereClase: false,
		},
	}
	for i, tp := range tiposProducto {
		if err := db.FirstOrCreate(
			&tiposProducto[i],
			TipoProducto{
				EmpresaID: tp.EmpresaID,
				Nombre:    tp.Nombre,
			},
		).Error; err != nil {
			fmt.Printf("Error creando TipoProducto %s: %v\n", tp.Nombre, err)
			return false
		}
	}

	areasDespacho := []AreaDespacho{
		{EmpresaID: 1, Nombre: "Cocina", ImprimirTicket: true, ImpresoraAsociada: "COCINA-PRINTER"},
		{EmpresaID: 1, Nombre: "Bar", ImprimirTicket: true, ImpresoraAsociada: "BAR-PRINTER"},
	}
	for i, area := range areasDespacho {
		if err := db.FirstOrCreate(&areasDespacho[i], AreaDespacho{EmpresaID: area.EmpresaID, Nombre: area.Nombre}).Error; err != nil {
			fmt.Printf("Error creando AreaDespacho %s: %v\n", area.Nombre, err)
			return false
		}
	}

	unidades := []Unidad{
		{EmpresaID: 1, Nombre: "Unidad", Simbolo: "UND", EsDecimal: false},
		{EmpresaID: 1, Nombre: "Litro", Simbolo: "L", EsDecimal: true},
	}
	for i, uni := range unidades {
		if err := db.FirstOrCreate(&unidades[i], Unidad{EmpresaID: uni.EmpresaID, Nombre: uni.Nombre, Simbolo: uni.Simbolo}).Error; err != nil {
			fmt.Printf("Error creando Unidad %s: %v\n", uni.Nombre, err)
			return false
		}
	}

	impuestos := []Impuesto{
		{EmpresaID: 1, NombreImpuesto: "IGV", Porcentaje: 18.0, IdeTributo: 1000, NomTributo: "IGV", CodTipTributo: "VAT"},
	}
	for i, imp := range impuestos {
		if err := db.FirstOrCreate(&impuestos[i], Impuesto{EmpresaID: imp.EmpresaID, NombreImpuesto: imp.NombreImpuesto}).Error; err != nil {
			fmt.Printf("Error creando Impuesto %s: %v\n", imp.NombreImpuesto, err)
			return false
		}
	}

	monitoresCocina := []MonitorCocina{
		{EmpresaID: 1, Monitor: "Cocina Principal"},
	}
	for i, mc := range monitoresCocina {
		if err := db.FirstOrCreate(&monitoresCocina[i], MonitorCocina{EmpresaID: mc.EmpresaID, Monitor: mc.Monitor}).Error; err != nil {
			fmt.Printf("Error creando MonitorCocina %s: %v\n", mc.Monitor, err)
			return false
		}
	}

	monitoresAuxiliares := []MonitorAuxiliar{
		{EmpresaID: 1, MonitorAuxiliar: "Caja"},
	}
	for i, ma := range monitoresAuxiliares {
		if err := db.FirstOrCreate(&monitoresAuxiliares[i], MonitorAuxiliar{EmpresaID: ma.EmpresaID, MonitorAuxiliar: ma.MonitorAuxiliar}).Error; err != nil {
			fmt.Printf("Error creando MonitorAuxiliar %s: %v\n", ma.MonitorAuxiliar, err)
			return false
		}
	}

	// 2. Sembrar entidades de categorías y clases
	categorias := []Categoria{
		{EmpresaID: 1, Nombre: "Platos Principales"},
		{EmpresaID: 1, Nombre: "Bebidas No Alcohólicas"},
	}
	for i, cat := range categorias {
		if err := db.FirstOrCreate(&categorias[i], Categoria{EmpresaID: cat.EmpresaID, Nombre: cat.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Categoria %s: %v\n", cat.Nombre, err)
			return false
		}
	}

	subCategorias := []SubCategoria{
		{CategoriaID: uint64(categorias[0].ID), Nombre: "General Platos", Activo: true},
		{CategoriaID: uint64(categorias[1].ID), Nombre: "General Bebidas", Activo: true},
	}
	for i, sc := range subCategorias {
		if err := db.FirstOrCreate(&subCategorias[i], SubCategoria{CategoriaID: sc.CategoriaID, Nombre: sc.Nombre}).Error; err != nil {
			fmt.Printf("Error creando SubCategoria %s: %v\n", sc.Nombre, err)
			return false
		}
	}

	clases := []Clase{
		{EmpresaID: 1, Nombre: "Default", Activo: true, OrdenVisualizacion: 1},
	}
	for i, cl := range clases {
		if err := db.FirstOrCreate(&clases[i], Clase{EmpresaID: cl.EmpresaID, Nombre: cl.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Clase %s: %v\n", cl.Nombre, err)
			return false
		}
	}

	// 3. Sembrar cartas y familias
	horaInicio := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
	horaFin := time.Date(2023, 1, 1, 22, 0, 0, 0, time.UTC)
	cartas := []Carta{
		{EmpresaID: 1, Nombre: "Carta Principal", Descripcion: "Carta principal", Imagen: "carta_principal.jpg", OrdenVisualizacion: 1, ControlHorarioActivo: false, HoraInicioServicio: horaInicio, HoraFinServicio: horaFin},
	}
	for i, carta := range cartas {
		if err := db.FirstOrCreate(&cartas[i], Carta{EmpresaID: carta.EmpresaID, Nombre: carta.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Carta %s: %v\n", carta.Nombre, err)
			return false
		}
	}

	familias := []Familia{
		{CartaID: uint64(cartas[0].ID), EmpresaID: 1, Nombre: "Platos Criollos", OrdenVisualizacion: 1, Imagen: "platos_criollos.jpg", Activo: true, Destacado: true},
	}
	for i, fam := range familias {
		if err := db.FirstOrCreate(&familias[i], Familia{CartaID: fam.CartaID, EmpresaID: fam.EmpresaID, Nombre: fam.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Familia %s: %v\n", fam.Nombre, err)
			return false
		}
	}

	// 4. Sembrar modificadores y opciones
	modificadores := []Modificador{
		{EmpresaID: 1, Nombre: "Nivel de Cocción", MinSelecciones: 1, MaxSelecciones: 1, MultipleSeleccion: false, Activo: true},
	}
	for i, mod := range modificadores {
		if err := db.FirstOrCreate(&modificadores[i], Modificador{EmpresaID: mod.EmpresaID, Nombre: mod.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Modificador %s: %v\n", mod.Nombre, err)
			return false
		}
	}

	var modCoccion Modificador
	db.Where("nombre = ? AND empresa_id = 1", "Nivel de Cocción").First(&modCoccion)

	opcionesModificador := []ModificadorOpcion{
		{ModificadorID: uint64(modCoccion.ID), Nombre: "Término Medio", PrecioAdicional: 0.0, OrdenVisualizacion: 1, Activo: true},
		{ModificadorID: uint64(modCoccion.ID), Nombre: "Tres Cuartos", PrecioAdicional: 0.0, OrdenVisualizacion: 2, Activo: true},
		{ModificadorID: uint64(modCoccion.ID), Nombre: "Bien Cocido", PrecioAdicional: 0.0, OrdenVisualizacion: 3, Activo: true},
	}
	for i, op := range opcionesModificador {
		if err := db.FirstOrCreate(&opcionesModificador[i], ModificadorOpcion{ModificadorID: op.ModificadorID, Nombre: op.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Opción de Modificador %s: %v\n", op.Nombre, err)
			return false
		}
	}

	// 5. Sembrar un Producto (usando SubCategoria en vez de Categoria)
	// Se recuperan IDs necesarios:
	var tipoPlatoID uint64
	db.Model(&TipoProducto{}).Where("nombre = ? AND empresa_id = 1", "Plato").Select("id").First(&tipoPlatoID)
	var areaCocinaID uint64
	db.Model(&AreaDespacho{}).Where("nombre = ? AND empresa_id = 1", "Cocina").Select("id").First(&areaCocinaID)
	var unidadUndID uint64
	db.Model(&Unidad{}).Where("nombre = ? AND empresa_id = 1", "Unidad").Select("id").First(&unidadUndID)
	var impuestoIGVID uint64
	db.Model(&Impuesto{}).Where("nombre_impuesto = ? AND empresa_id = 1", "IGV").Select("id").First(&impuestoIGVID)
	var monitorCocinaID uint64
	db.Model(&MonitorCocina{}).Where("monitor = ? AND empresa_id = 1", "Cocina Principal").Select("id").First(&monitorCocinaID)
	var monitorAuxiliarID uint64
	db.Model(&MonitorAuxiliar{}).Where("monitor_auxiliar = ? AND empresa_id = 1", "Caja").Select("id").First(&monitorAuxiliarID)
	var subCategoriaPlatoID uint64
	db.Model(&SubCategoria{}).Where("nombre = ? AND categoria_id = ?", "General Platos", categorias[0].ID).Select("id").First(&subCategoriaPlatoID)
	var claseDefaultID uint64
	db.Model(&Clase{}).Where("nombre = ? AND empresa_id = 1", "Default").Select("id").First(&claseDefaultID)

	productos := []Producto{
		{
			EmpresaID:         1,
			UsuarioID:         1,
			TipoProductoID:    tipoPlatoID,
			SubCategoriaID:    &subCategoriaPlatoID,
			AreaDespachoID:    &areaCocinaID,
			ClaseID:           &claseDefaultID,
			UnidadControlID:   unidadUndID,
			ImpuestoID:        impuestoIGVID,
			MonitorCocinaID:   &monitorCocinaID,
			MonitorAuxiliarID: &monitorAuxiliarID,
			CodigoSunatID:     1,
			Nombre:            "Lomo Saltado",
			Descripcion:       "Lomo saltado tradicional",
			Codigo:            "P001",
			CodigoBarras:      "0001",
			TiempoPreparacion: 15,
			Precio:            35.90,
		},
		// Productos adicionales para pruebas de búsqueda
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Arroz Chaufa", Descripcion: "Arroz chaufa de pollo", Codigo: "P002", CodigoBarras: "0002", TiempoPreparacion: 12, Precio: 28.50,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Ceviche Mixto", Descripcion: "Ceviche de pescado y mariscos", Codigo: "P003", CodigoBarras: "0003", TiempoPreparacion: 10, Precio: 32.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Pollo a la Brasa", Descripcion: "Pollo a la brasa con papas", Codigo: "P004", CodigoBarras: "0004", TiempoPreparacion: 25, Precio: 40.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Tallarines Verdes", Descripcion: "Tallarines con salsa de albahaca", Codigo: "P005", CodigoBarras: "0005", TiempoPreparacion: 14, Precio: 27.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Ají de Gallina", Descripcion: "Ají de gallina tradicional", Codigo: "P006", CodigoBarras: "0006", TiempoPreparacion: 13, Precio: 26.50,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Papa a la Huancaína", Descripcion: "Entrada de papa a la huancaína", Codigo: "P007", CodigoBarras: "0007", TiempoPreparacion: 8, Precio: 15.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Causa Limeña", Descripcion: "Causa de pollo", Codigo: "P008", CodigoBarras: "0008", TiempoPreparacion: 9, Precio: 16.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Anticuchos", Descripcion: "Anticuchos de corazón", Codigo: "P009", CodigoBarras: "0009", TiempoPreparacion: 11, Precio: 18.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Seco de Cabrito", Descripcion: "Seco de cabrito con frejoles", Codigo: "P010", CodigoBarras: "0010", TiempoPreparacion: 18, Precio: 38.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Tacu Tacu", Descripcion: "Tacu tacu con lomo", Codigo: "P011", CodigoBarras: "0011", TiempoPreparacion: 16, Precio: 30.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Chicharrón de Pescado", Descripcion: "Chicharrón de pescado crocante", Codigo: "P012", CodigoBarras: "0012", TiempoPreparacion: 12, Precio: 22.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Jalea Mixta", Descripcion: "Jalea de mariscos", Codigo: "P013", CodigoBarras: "0013", TiempoPreparacion: 15, Precio: 36.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Ensalada Rusa", Descripcion: "Ensalada de verduras y mayonesa", Codigo: "P014", CodigoBarras: "0014", TiempoPreparacion: 7, Precio: 12.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Sopa Criolla", Descripcion: "Sopa criolla peruana", Codigo: "P015", CodigoBarras: "0015", TiempoPreparacion: 10, Precio: 17.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Chaufa de Mariscos", Descripcion: "Arroz chaufa con mariscos", Codigo: "P016", CodigoBarras: "0016", TiempoPreparacion: 13, Precio: 33.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Sudado de Pescado", Descripcion: "Sudado de pescado fresco", Codigo: "P017", CodigoBarras: "0017", TiempoPreparacion: 14, Precio: 29.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Milanesa de Pollo", Descripcion: "Milanesa de pollo crocante", Codigo: "P018", CodigoBarras: "0018", TiempoPreparacion: 12, Precio: 25.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Bistec a lo Pobre", Descripcion: "Bistec con huevo y plátano", Codigo: "P019", CodigoBarras: "0019", TiempoPreparacion: 15, Precio: 34.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Chicha Morada", Descripcion: "Bebida tradicional peruana", Codigo: "B001", CodigoBarras: "1001", TiempoPreparacion: 2, Precio: 6.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Inca Kola", Descripcion: "Bebida gaseosa peruana", Codigo: "B002", CodigoBarras: "1002", TiempoPreparacion: 1, Precio: 7.00,
		},
		{
			EmpresaID: 1, UsuarioID: 1, TipoProductoID: tipoPlatoID, SubCategoriaID: &subCategoriaPlatoID, AreaDespachoID: &areaCocinaID, ClaseID: &claseDefaultID, UnidadControlID: unidadUndID, ImpuestoID: impuestoIGVID, MonitorCocinaID: &monitorCocinaID, MonitorAuxiliarID: &monitorAuxiliarID, CodigoSunatID: 1, Nombre: "Agua Mineral", Descripcion: "Agua embotellada", Codigo: "B003", CodigoBarras: "1003", TiempoPreparacion: 1, Precio: 5.00,
		},
	}
	for i, p := range productos {
		if err := db.FirstOrCreate(&productos[i], Producto{EmpresaID: p.EmpresaID, Codigo: p.Codigo}).Error; err != nil {
			fmt.Printf("Error creando Producto %s: %v\n", p.Nombre, err)
			return false
		}
	}

	// 6. Sembrar Presentaciones para el Producto
	var productoID uint64 = uint64(productos[0].ID)
	presentaciones := []Presentacion{
		{
			ProductoID:            productoID,
			Nombre:                "Lomo Saltado Regular",
			Descripcion:           "Porción regular",
			CantidadUnidadControl: 1,
			Precio:                35.90,
			EsPresentacionBase:    true,
		},
		{
			ProductoID:            productoID,
			Nombre:                "Lomo Saltado Grande",
			Descripcion:           "Porción grande",
			CantidadUnidadControl: 1.5,
			Precio:                45.90,
			EsPresentacionBase:    false,
		},
	}
	for i, pres := range presentaciones {
		if err := db.FirstOrCreate(&presentaciones[i], Presentacion{ProductoID: pres.ProductoID, Nombre: pres.Nombre}).Error; err != nil {
			fmt.Printf("Error creando Presentacion %s: %v\n", pres.Nombre, err)
			return false
		}
	}

	// 7. Sembrar Listas de Precios y su relación con Presentaciones
	listasPrecios := []ListaPrecios{
		{EmpresaID: 1, Nombre: "Lista Estándar", Descripcion: "Precios estándar", EsPrecioPorMayor: false, EsPrecioEspecial: false},
		{EmpresaID: 1, Nombre: "Lista Mayorista", Descripcion: "Precios con descuento", EsPrecioPorMayor: true, EsPrecioEspecial: false},
	}
	for i, lp := range listasPrecios {
		if err := db.FirstOrCreate(&listasPrecios[i], ListaPrecios{EmpresaID: lp.EmpresaID, Nombre: lp.Nombre}).Error; err != nil {
			fmt.Printf("Error creando ListaPrecios %s: %v\n", lp.Nombre, err)
			return false
		}
	}

	// Relacionar cada Presentacion con cada Lista de Precios
	for _, pres := range presentaciones {
		// Precio para Lista Estándar
		lpEst := ListaPreciosPresentacion{
			ListaPreciosID: uint64(listasPrecios[0].ID),
			PresentacionID: uint64(pres.ID),
			Precio:         pres.Precio,
		}
		if result := db.FirstOrCreate(&lpEst, ListaPreciosPresentacion{ListaPreciosID: lpEst.ListaPreciosID, PresentacionID: lpEst.PresentacionID}); result.Error != nil {
			fmt.Printf("Error creando relación ListaPreciosPresentacion: %v\n", result.Error)
			return false
		}

		// Precio para Lista Mayorista (10% descuento)
		lpMayor := ListaPreciosPresentacion{
			ListaPreciosID: uint64(listasPrecios[1].ID),
			PresentacionID: uint64(pres.ID),
			Precio:         pres.Precio * 0.9,
		}
		if result := db.FirstOrCreate(&lpMayor, ListaPreciosPresentacion{ListaPreciosID: lpMayor.ListaPreciosID, PresentacionID: lpMayor.PresentacionID}); result.Error != nil {
			fmt.Printf("Error creando relación ListaPreciosPresentacion mayorista: %v\n", result.Error)
			return false
		}
	}

	// 8. Sembrar otros datos relacionados al producto
	productoMedias := []ProductoMedia{
		{ProductoID: productoID, URL: "lomo_saltado.jpg", Tipo: "IMAGEN", Orden: 1, EsPrincipal: true},
	}
	for i, pm := range productoMedias {
		if result := db.FirstOrCreate(&productoMedias[i], ProductoMedia{ProductoID: pm.ProductoID, URL: pm.URL}); result.Error != nil {
			fmt.Printf("Error creando ProductoMedia: %v\n", result.Error)
			return false
		}
	}

	comentariosProducto := []ComentarioProducto{
		{ProductoID: &productoID, Comentario: "Sin cebolla", EsComentarioGenerico: false},
		{ProductoID: nil, Comentario: "Para llevar", EsComentarioGenerico: true},
	}
	for i, cp := range comentariosProducto {
		var query ComentarioProducto
		if cp.ProductoID != nil {
			query = ComentarioProducto{ProductoID: cp.ProductoID, Comentario: cp.Comentario}
		} else {
			query = ComentarioProducto{EsComentarioGenerico: true, Comentario: cp.Comentario}
		}
		if result := db.FirstOrCreate(&comentariosProducto[i], query); result.Error != nil {
			fmt.Printf("Error creando ComentarioProducto: %v\n", result.Error)
			return false
		}
	}

	var presRegularID uint64
	db.Model(&Presentacion{}).Where("nombre = ? AND producto_id = ?", "Lomo Saltado Regular", productoID).Select("id").First(&presRegularID)

	precioHistorico := PrecioHistorico{
		PresentacionID: presRegularID,
		PrecioAnterior: 32.90,
		PrecioNuevo:    35.90,
		UsuarioID:      1,
		PrecioActivo:   true,
	}
	if result := db.FirstOrCreate(&precioHistorico, PrecioHistorico{PresentacionID: precioHistorico.PresentacionID, PrecioNuevo: precioHistorico.PrecioNuevo}); result.Error != nil {
		fmt.Printf("Error creando PrecioHistorico: %v\n", result.Error)
		return false
	}

	preciosDescuento := []PrecioDescuentoPorCantidad{
		{ProductoID: productoID, CantidadMinima: 3, CantidadMaxima: 10, UsaPrecio: true, PrecioAlternativo: 33.90, Descuento: 2.00},
	}
	for i, pd := range preciosDescuento {
		if result := db.FirstOrCreate(&preciosDescuento[i], PrecioDescuentoPorCantidad{ProductoID: pd.ProductoID, CantidadMinima: pd.CantidadMinima}); result.Error != nil {
			fmt.Printf("Error creando PrecioDescuentoPorCantidad: %v\n", result.Error)
			return false
		}
	}

	// 9. Sembrar Grupo de Conteo y su relación con el Producto
	grupoConteo := GrupoConteo{
		EmpresaID:   1,
		Nombre:      "Inventario Mensual",
		Descripcion: "Conteo mensual de inventario",
		AlmacenID:   1,
		OrdenConteo: 1,
	}
	if result := db.FirstOrCreate(&grupoConteo, GrupoConteo{EmpresaID: grupoConteo.EmpresaID, Nombre: grupoConteo.Nombre}); result.Error != nil {
		fmt.Printf("Error creando GrupoConteo: %v\n", result.Error)
		return false
	}

	productoGrupo := ProductoGrupoConteo{
		ProductoID:    productoID,
		GrupoConteoID: uint64(grupoConteo.ID),
	}
	if result := db.FirstOrCreate(&productoGrupo, ProductoGrupoConteo{ProductoID: productoGrupo.ProductoID, GrupoConteoID: productoGrupo.GrupoConteoID}); result.Error != nil {
		fmt.Printf("Error creando ProductoGrupoConteo: %v\n", result.Error)
		return false
	}

	// 10. Sembrar PresentacionPeso, PresentacionTiempo y PresentacionCompuesta
	peso := PresentacionPeso{
		PresentacionID: presRegularID,
		PrecioXKilo:    8.60,
		PesoMinimo:     1.0,
		PesoMaximo:     1.5,
		Activo:         true,
	}
	if result := db.FirstOrCreate(&peso, PresentacionPeso{PresentacionID: peso.PresentacionID}); result.Error != nil {
		fmt.Printf("Error creando PresentacionPeso: %v\n", result.Error)
		return false
	}

	tiempo := PresentacionTiempo{
		PresentacionID: presRegularID,
		PrecioXHora:    10.0,
		TiempoMinimo:   1.0,
		TiempoMaximo:   3.0,
		FraccionCobro:  0.5,
		Activo:         true,
	}
	if result := db.FirstOrCreate(&tiempo, PresentacionTiempo{PresentacionID: tiempo.PresentacionID}); result.Error != nil {
		fmt.Printf("Error creando PresentacionTiempo: %v\n", result.Error)
		return false
	}

	// Ejemplo de PresentacionCompuesta (si existiera un componente)
	// Aquí se omite si no hay otro producto para componer, se puede expandir según necesidades

	// 11. Sembrar ProductoModificador
	prodMod := ProductoModificador{
		ProductoID:         productoID,
		Modificador:        &modCoccion,
		Required:           true,
		OrdenVisualizacion: 1,
		Activo:             true,
		Deleted:            false,
	}
	if result := db.FirstOrCreate(&prodMod, ProductoModificador{ProductoID: prodMod.ProductoID, ModificadorID: uint64(modCoccion.ID)}); result.Error != nil {
		fmt.Printf("Error creando ProductoModificador: %v\n", result.Error)
		return false
	}

	// 12. Sembrar Oferta, GrupoOferta y OfertaPresentacion
	oferta := Oferta{
		EmpresaID:      1,
		PresentacionID: presRegularID,
		FechaInicio:    time.Now().Add(-24 * time.Hour),
		FechaFin:       time.Now().Add(7 * 24 * time.Hour),
		StockOfertado:  100,
		StockActual:    100,
		Activo:         true,
	}
	if result := db.FirstOrCreate(&oferta, Oferta{EmpresaID: oferta.EmpresaID, PresentacionID: oferta.PresentacionID}); result.Error != nil {
		fmt.Printf("Error creando Oferta: %v\n", result.Error)
		return false
	}

	var ofertaID uint64 = uint64(oferta.ID)

	grupoOferta := GrupoOferta{
		OfertaID:           ofertaID,
		Nombre:             "Promo Especial",
		MinSelecciones:     1,
		MaxSelecciones:     2,
		PrecioBase:         7.00,
		OrdenVisualizacion: 1,
		Activo:             true,
	}
	if result := db.FirstOrCreate(&grupoOferta, GrupoOferta{OfertaID: grupoOferta.OfertaID, Nombre: grupoOferta.Nombre}); result.Error != nil {
		fmt.Printf("Error creando GrupoOferta: %v\n", result.Error)
		return false
	}

	ofertaPres := OfertaPresentacion{
		OfertaID:        ofertaID,
		GrupoOfertaID:   uint64(grupoOferta.ID),
		PresentacionID:  presRegularID,
		PrecioAdicional: 0.0,
		Activo:          true,
		Deleted:         false,
	}
	if result := db.FirstOrCreate(&ofertaPres, OfertaPresentacion{OfertaID: ofertaPres.OfertaID, GrupoOfertaID: ofertaPres.GrupoOfertaID, PresentacionID: ofertaPres.PresentacionID}); result.Error != nil {
		fmt.Printf("Error creando OfertaPresentacion: %v\n", result.Error)
		return false
	}

	fmt.Println("¡Datos semilla (v2) insertados exitosamente!")
	return true
}

// SeedersHandler expone vía HTTP la ejecución de Seeders.
func SeedersHandler(c *fiber.Ctx) error {
	if !Seeders() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":    "error",
			"message":   "Error al cargar seeders v2",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Seeders v2 ejecutados exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// DropTablesHandler maneja la eliminación de todas las tablas via HTTP.
func DropTablesHandler(c *fiber.Ctx) error {
	database.ConnectDatabase()
	db := database.DBconn
	if err := DropTables(db); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":    "error",
			"message":   "Error al eliminar tablas: " + err.Error(),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "success",
		"message":   "Tablas eliminadas exitosamente",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// RegisterSeedersV2Route registra la ruta para Seeders.
func RegisterSeedersV2Route(app *fiber.App) {
	app.Get("/api/v1/producto/load-seeders", SeedersHandler)
	app.Get("/api/v1/producto/drop-tables", DropTablesHandler)
}

func init() {
	registry.RegisterModule(RegisterSeedersV2Route)
}
