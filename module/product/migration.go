package product

import (
	"practicev2/registry"

	"gorm.io/gorm"
)

// Migrate ejecuta la migración de las tablas de productos en el orden correcto.
func Migrate(db *gorm.DB) error {
	// 1. Tablas base sin dependencias
	if err := db.AutoMigrate(&TipoProducto{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&AreaDespacho{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Unidad{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&MonitorCocina{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&MonitorAuxiliar{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Impuesto{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Categoria{}); err != nil {
		return err
	}

	// 2. Tablas con dependencias simples
	if err := db.AutoMigrate(&SubCategoria{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Clase{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Modificador{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ModificadorOpcion{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ListaPrecios{}); err != nil {
		return err
	}

	// 3. Tablas con múltiples dependencias o dependencias en cascada
	if err := db.AutoMigrate(&Producto{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Presentacion{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ComentarioProducto{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&GrupoConteo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Carta{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Familia{}); err != nil {
		return err
	}

	// 4. Tablas de relaciones y detalles
	if err := db.AutoMigrate(&ProductoGrupoConteo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ProductoFamilia{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ProductoMedia{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&PrecioHistorico{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&PrecioDescuentoPorCantidad{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ProductoModificador{}); err != nil {
		return err
	}

	// 5. Tablas relacionadas con ofertas
	if err := db.AutoMigrate(&Oferta{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&GrupoOferta{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&OfertaPresentacion{}); err != nil {
		return err
	}

	// 6. Tablas de configuración de presentaciones
	if err := db.AutoMigrate(&PresentacionCompuesta{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&PresentacionPeso{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&PresentacionTiempo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&ListaPreciosPresentacion{}); err != nil {
		return err
	}

	return nil
}

// DropTables elimina todas las tablas siguiendo el orden correcto según sus relaciones.
func DropTables(db *gorm.DB) error {
	// 6. Eliminar primero las tablas de configuración de presentaciones
	if err := db.Migrator().DropTable(&ListaPreciosPresentacion{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&PresentacionTiempo{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&PresentacionPeso{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&PresentacionCompuesta{}); err != nil {
		return err
	}

	// 5. Eliminar tablas relacionadas con ofertas
	if err := db.Migrator().DropTable(&OfertaPresentacion{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&GrupoOferta{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Oferta{}); err != nil {
		return err
	}

	// 4. Eliminar tablas de relaciones y detalles
	if err := db.Migrator().DropTable(&ProductoModificador{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&PrecioDescuentoPorCantidad{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&PrecioHistorico{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&ProductoMedia{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&ProductoFamilia{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&ProductoGrupoConteo{}); err != nil {
		return err
	}

	// 3. Eliminar tablas con múltiples dependencias
	if err := db.Migrator().DropTable(&Familia{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Carta{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&GrupoConteo{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&ComentarioProducto{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Presentacion{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Producto{}); err != nil {
		return err
	}

	// 2. Eliminar tablas con dependencias simples
	if err := db.Migrator().DropTable(&ListaPrecios{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&ModificadorOpcion{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Modificador{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Clase{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&SubCategoria{}); err != nil {
		return err
	}

	// 1. Eliminar tablas base sin dependencias
	if err := db.Migrator().DropTable(&Categoria{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Impuesto{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&MonitorAuxiliar{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&MonitorCocina{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&Unidad{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&AreaDespacho{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&TipoProducto{}); err != nil {
		return err
	}

	return nil
}

func init() {
	// Se registra la migración con orden 3 (posterior a usuarios).
	registry.RegisterMigration(3, Migrate)
}
