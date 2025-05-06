package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ProductoMediaRepository define la interfaz para operaciones de persistencia de ProductoMedia
type ProductoMediaRepository interface {
	GetAllProductoMedias(empresaID uint64) ([]ProductoMedia, error)
	GetProductoMediasByProductoID(productoID uint64, empresaID uint64) ([]ProductoMedia, error)
	GetProductoMedia(id uint64, empresaID uint64) (ProductoMedia, error)
	CreateProductoMedia(productoMedia *ProductoMedia) error
	UpdateProductoMedia(id uint64, empresaID uint64, productoMedia *ProductoMedia) (ProductoMedia, error)
	DeleteProductoMedia(id uint64, empresaID uint64) error
	SeedProductoMedias(productoMedias []ProductoMedia) error
	PatchProductoMedia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoMedia, error)
	SetMainImage(id uint64, productoID uint64, empresaID uint64) error
	ReorderMedia(productoID uint64, empresaID uint64, mediaIDs []uint64) error
}

type productoMediaRepo struct {
	DB *gorm.DB
}

// NewProductoMediaRepository crea una instancia del repositorio
func NewProductoMediaRepository() ProductoMediaRepository {
	return &productoMediaRepo{DB: database.DBconn}
}

// GetAllProductoMedias obtiene todos los medios asociados a productos para una empresa
func (r *productoMediaRepo) GetAllProductoMedias(empresaID uint64) ([]ProductoMedia, error) {
	var productoMedias []ProductoMedia
	err := r.DB.
		Joins("JOIN productos ON producto_media.producto_id = productos.id").
		Where("productos.empresa_id = ?", empresaID).
		Order("producto_media.producto_id, producto_media.orden").
		Find(&productoMedias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar medios de productos: %w", err)
	}
	return productoMedias, nil
}

// GetProductoMediasByProductoID obtiene todos los medios para un producto específico
func (r *productoMediaRepo) GetProductoMediasByProductoID(productoID uint64, empresaID uint64) ([]ProductoMedia, error) {
	var productoMedias []ProductoMedia
	err := r.DB.
		Joins("JOIN productos ON producto_media.producto_id = productos.id").
		Where("producto_media.producto_id = ? AND productos.empresa_id = ?", productoID, empresaID).
		Order("producto_media.orden").
		Find(&productoMedias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar medios para producto ID %d: %w", productoID, err)
	}
	return productoMedias, nil
}

// GetProductoMedia obtiene un medio específico
func (r *productoMediaRepo) GetProductoMedia(id uint64, empresaID uint64) (ProductoMedia, error) {
	var productoMedia ProductoMedia
	err := r.DB.
		Joins("JOIN productos ON producto_media.producto_id = productos.id").
		Where("producto_media.id = ? AND productos.empresa_id = ?", id, empresaID).
		First(&productoMedia).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoMedia, fmt.Errorf("medio con id %d no encontrado", id)
		}
		return productoMedia, fmt.Errorf("error al recuperar medio con id %d: %w", id, err)
	}
	return productoMedia, nil
}

// CreateProductoMedia crea un nuevo medio para un producto
func (r *productoMediaRepo) CreateProductoMedia(productoMedia *ProductoMedia) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el producto existe y pertenece a la empresa
		var producto Producto
		if err := tx.Where("id = ?", productoMedia.ProductoID).First(&producto).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("producto con id %d no encontrado", productoMedia.ProductoID)
			}
			return fmt.Errorf("error al verificar producto: %w", err)
		}

		// Si es la primera imagen o se marca como principal, asegurar que no haya otra principal
		if productoMedia.EsPrincipal {
			if err := tx.Model(&ProductoMedia{}).
				Where("producto_id = ? AND es_principal = ?", productoMedia.ProductoID, true).
				Update("es_principal", false).Error; err != nil {
				return fmt.Errorf("error al actualizar estado de imagen principal: %w", err)
			}
		}

		// Determinar el siguiente orden para el medio
		var maxOrden int
		if err := tx.Model(&ProductoMedia{}).
			Where("producto_id = ?", productoMedia.ProductoID).
			Select("COALESCE(MAX(orden), 0)").
			Row().
			Scan(&maxOrden); err != nil {
			return fmt.Errorf("error al determinar orden: %w", err)
		}
		productoMedia.Orden = maxOrden + 1

		// Crear el medio
		return tx.Create(productoMedia).Error
	})
	if err != nil {
		return fmt.Errorf("error al crear medio para producto: %w", err)
	}
	return nil
}

// UpdateProductoMedia actualiza un medio existente
func (r *productoMediaRepo) UpdateProductoMedia(id uint64, empresaID uint64, updatedProductoMedia *ProductoMedia) (ProductoMedia, error) {
	var productoMedia ProductoMedia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Recuperar el medio actual para verificar pertenencia
		if err := tx.
			Joins("JOIN productos ON producto_media.producto_id = productos.id").
			Where("producto_media.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&productoMedia).Error; err != nil {
			return err
		}

		// Si se cambia a principal, actualizar otros medios
		if updatedProductoMedia.EsPrincipal && !productoMedia.EsPrincipal {
			if err := tx.Model(&ProductoMedia{}).
				Where("producto_id = ? AND es_principal = ? AND id != ?", productoMedia.ProductoID, true, id).
				Update("es_principal", false).Error; err != nil {
				return fmt.Errorf("error al actualizar estado de imagen principal: %w", err)
			}
		}

		return tx.Model(&productoMedia).Updates(updatedProductoMedia).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoMedia, fmt.Errorf("medio con id %d no encontrado", id)
		}
		return productoMedia, fmt.Errorf("error al actualizar medio con id %d: %w", id, err)
	}
	return productoMedia, nil
}

// DeleteProductoMedia elimina un medio
func (r *productoMediaRepo) DeleteProductoMedia(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var productoMedia ProductoMedia
		if err := tx.
			Joins("JOIN productos ON producto_media.producto_id = productos.id").
			Where("producto_media.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&productoMedia).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("medio con id %d no encontrado", id)
			}
			return fmt.Errorf("error al buscar medio: %w", err)
		}

		// Si se elimina un medio principal, establecer otro como principal si hay más medios
		if productoMedia.EsPrincipal {
			var count int64
			if err := tx.Model(&ProductoMedia{}).
				Where("producto_id = ? AND id != ?", productoMedia.ProductoID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("error al contar medios: %w", err)
			}

			if count > 0 {
				// Hay otros medios, establecer el primero como principal
				var otroMedia ProductoMedia
				if err := tx.Where("producto_id = ? AND id != ?", productoMedia.ProductoID, id).
					Order("orden").
					First(&otroMedia).Error; err != nil {
					return fmt.Errorf("error al buscar otro medio: %w", err)
				}

				if err := tx.Model(&otroMedia).Update("es_principal", true).Error; err != nil {
					return fmt.Errorf("error al actualizar nuevo medio principal: %w", err)
				}
			}
		}

		// Eliminar el medio
		return tx.Delete(&productoMedia).Error
	})
	if err != nil {
		return fmt.Errorf("error al eliminar medio: %w", err)
	}
	return nil
}

// SeedProductoMedias inicializa múltiples medios para productos
func (r *productoMediaRepo) SeedProductoMedias(productoMedias []ProductoMedia) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&productoMedias).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar medios de productos: %w", err)
	}
	return nil
}

// PatchProductoMedia actualiza parcialmente un medio
func (r *productoMediaRepo) PatchProductoMedia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoMedia, error) {
	var productoMedia ProductoMedia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Recuperar el medio actual para verificar pertenencia
		if err := tx.
			Joins("JOIN productos ON producto_media.producto_id = productos.id").
			Where("producto_media.id = ? AND productos.empresa_id = ?", id, empresaID).
			First(&productoMedia).Error; err != nil {
			return err
		}

		// Si se cambia a principal, actualizar otros medios
		if esPrincipal, ok := fields["es_principal"].(bool); ok && esPrincipal && !productoMedia.EsPrincipal {
			if err := tx.Model(&ProductoMedia{}).
				Where("producto_id = ? AND es_principal = ? AND id != ?", productoMedia.ProductoID, true, id).
				Update("es_principal", false).Error; err != nil {
				return fmt.Errorf("error al actualizar estado de imagen principal: %w", err)
			}
		}

		return tx.Model(&productoMedia).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoMedia, fmt.Errorf("medio con id %d no encontrado", id)
		}
		return productoMedia, fmt.Errorf("error al actualizar parcialmente medio con id %d: %w", id, err)
	}
	return productoMedia, nil
}

// SetMainImage establece una imagen como principal para un producto
func (r *productoMediaRepo) SetMainImage(id uint64, productoID uint64, empresaID uint64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el medio existe y pertenece al producto
		var count int64
		if err := tx.
			Joins("JOIN productos ON producto_media.producto_id = productos.id").
			Where("producto_media.id = ? AND producto_media.producto_id = ? AND productos.empresa_id = ?", id, productoID, empresaID).
			Count(&count).Error; err != nil {
			return fmt.Errorf("error al verificar medio: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("medio con id %d no encontrado o no pertenece al producto %d", id, productoID)
		}

		// Quitar el atributo principal de todos los medios del producto
		if err := tx.Model(&ProductoMedia{}).
			Where("producto_id = ?", productoID).
			Update("es_principal", false).Error; err != nil {
			return fmt.Errorf("error al actualizar estado de imagen principal: %w", err)
		}

		// Establecer el medio seleccionado como principal
		if err := tx.Model(&ProductoMedia{}).
			Where("id = ?", id).
			Update("es_principal", true).Error; err != nil {
			return fmt.Errorf("error al establecer imagen principal: %w", err)
		}

		return nil
	})
}

// ReorderMedia reordena los medios de un producto
func (r *productoMediaRepo) ReorderMedia(productoID uint64, empresaID uint64, mediaIDs []uint64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que el producto pertenece a la empresa
		var producto Producto
		if err := tx.Where("id = ? AND empresa_id = ?", productoID, empresaID).First(&producto).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("producto con id %d no encontrado o no pertenece a la empresa", productoID)
			}
			return fmt.Errorf("error al verificar producto: %w", err)
		}

		// Verificar que todos los medios pertenecen al producto
		for i, mediaID := range mediaIDs {
			var media ProductoMedia
			if err := tx.Where("id = ? AND producto_id = ?", mediaID, productoID).First(&media).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("medio con id %d no encontrado o no pertenece al producto %d", mediaID, productoID)
				}
				return fmt.Errorf("error al verificar medio: %w", err)
			}

			// Actualizar el orden
			if err := tx.Model(&media).Update("orden", i+1).Error; err != nil {
				return fmt.Errorf("error al actualizar orden del medio %d: %w", mediaID, err)
			}
		}

		return nil
	})
}
