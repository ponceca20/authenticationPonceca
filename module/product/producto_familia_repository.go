package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ProductoFamiliaRepository define la interfaz para operaciones de persistencia de ProductoFamilia
type ProductoFamiliaRepository interface {
	GetAllProductoFamilias(empresaID uint64) ([]ProductoFamilia, error)
	GetProductoFamiliasByFamiliaID(familiaID uint64, empresaID uint64) ([]ProductoFamilia, error)
	GetProductoFamiliasByProductoID(productoID uint64, empresaID uint64) ([]ProductoFamilia, error)
	GetProductoFamilia(id uint64, empresaID uint64) (ProductoFamilia, error)
	CreateProductoFamilia(productoFamilia *ProductoFamilia) error
	UpdateProductoFamilia(id uint64, empresaID uint64, productoFamilia *ProductoFamilia) (ProductoFamilia, error)
	DeleteProductoFamilia(id uint64, empresaID uint64) error
	SeedProductoFamilias(productoFamilias []ProductoFamilia) error
	PatchProductoFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoFamilia, error)
}

type productoFamiliaRepo struct {
	DB *gorm.DB
}

// NewProductoFamiliaRepository crea una instancia del repositorio
func NewProductoFamiliaRepository() ProductoFamiliaRepository {
	return &productoFamiliaRepo{DB: database.DBconn}
}

// GetAllProductoFamilias obtiene todas las relaciones ProductoFamilia para una empresa
func (r *productoFamiliaRepo) GetAllProductoFamilias(empresaID uint64) ([]ProductoFamilia, error) {
	var productoFamilias []ProductoFamilia
	err := r.DB.
		Preload("Producto").
		Preload("Familia").
		Joins("JOIN familias ON producto_familias.familia_id = familias.id").
		Where("familias.empresa_id = ?", empresaID).
		Find(&productoFamilias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar relaciones producto-familia: %w", err)
	}
	return productoFamilias, nil
}

// GetProductoFamiliasByFamiliaID obtiene todas las relaciones ProductoFamilia para una familia específica
func (r *productoFamiliaRepo) GetProductoFamiliasByFamiliaID(familiaID uint64, empresaID uint64) ([]ProductoFamilia, error) {
	var productoFamilias []ProductoFamilia
	err := r.DB.
		Preload("Producto").
		Preload("Familia").
		Joins("JOIN familias ON producto_familias.familia_id = familias.id").
		Where("producto_familias.familia_id = ? AND familias.empresa_id = ?", familiaID, empresaID).
		Find(&productoFamilias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar relaciones producto-familia para familia ID %d: %w", familiaID, err)
	}
	return productoFamilias, nil
}

// GetProductoFamiliasByProductoID obtiene todas las relaciones ProductoFamilia para un producto específico
func (r *productoFamiliaRepo) GetProductoFamiliasByProductoID(productoID uint64, empresaID uint64) ([]ProductoFamilia, error) {
	var productoFamilias []ProductoFamilia
	err := r.DB.
		Preload("Producto").
		Preload("Familia").
		Joins("JOIN productos ON producto_familias.producto_id = productos.id").
		Where("producto_familias.producto_id = ? AND productos.empresa_id = ?", productoID, empresaID).
		Find(&productoFamilias).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar relaciones producto-familia para producto ID %d: %w", productoID, err)
	}
	return productoFamilias, nil
}

// GetProductoFamilia obtiene una relación ProductoFamilia específica
func (r *productoFamiliaRepo) GetProductoFamilia(id uint64, empresaID uint64) (ProductoFamilia, error) {
	var productoFamilia ProductoFamilia
	err := r.DB.
		Preload("Producto").
		Preload("Familia").
		Joins("JOIN familias ON producto_familias.familia_id = familias.id").
		Where("producto_familias.id = ? AND familias.empresa_id = ?", id, empresaID).
		First(&productoFamilia).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoFamilia, fmt.Errorf("relación producto-familia con id %d no encontrada", id)
		}
		return productoFamilia, fmt.Errorf("error al recuperar relación producto-familia con id %d: %w", id, err)
	}
	return productoFamilia, nil
}

// CreateProductoFamilia crea una nueva relación ProductoFamilia
func (r *productoFamiliaRepo) CreateProductoFamilia(productoFamilia *ProductoFamilia) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar que la familia existe y pertenece a la empresa
		var familia Familia
		if err := tx.Where("id = ?", productoFamilia.FamiliaID).First(&familia).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("familia con id %d no encontrada", productoFamilia.FamiliaID)
			}
			return fmt.Errorf("error al verificar familia: %w", err)
		}

		// Verificar que el producto existe
		var producto Producto
		if err := tx.Where("id = ?", productoFamilia.ProductoID).First(&producto).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("producto con id %d no encontrado", productoFamilia.ProductoID)
			}
			return fmt.Errorf("error al verificar producto: %w", err)
		}

		// Verificar que el producto y la familia pertenecen a la misma empresa
		if producto.EmpresaID != familia.EmpresaID {
			return fmt.Errorf("el producto y la familia deben pertenecer a la misma empresa")
		}

		// Verificar que no exista ya una relación entre este producto y familia
		var count int64
		if err := tx.Model(&ProductoFamilia{}).Where("producto_id = ? AND familia_id = ?", productoFamilia.ProductoID, productoFamilia.FamiliaID).Count(&count).Error; err != nil {
			return fmt.Errorf("error al verificar duplicidad: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("ya existe una relación entre el producto %d y la familia %d", productoFamilia.ProductoID, productoFamilia.FamiliaID)
		}

		return tx.Create(productoFamilia).Error
	})
	if err != nil {
		return fmt.Errorf("error al crear relación producto-familia: %w", err)
	}
	return nil
}

// UpdateProductoFamilia actualiza una relación ProductoFamilia existente
func (r *productoFamiliaRepo) UpdateProductoFamilia(id uint64, empresaID uint64, updatedProductoFamilia *ProductoFamilia) (ProductoFamilia, error) {
	var productoFamilia ProductoFamilia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Primero recuperamos la relación actual para verificar que pertenece a la empresa
		if err := tx.
			Joins("JOIN familias ON producto_familias.familia_id = familias.id").
			Where("producto_familias.id = ? AND familias.empresa_id = ?", id, empresaID).
			First(&productoFamilia).Error; err != nil {
			return err
		}

		// Si se está cambiando la familia o el producto, hacer verificaciones adicionales
		if updatedProductoFamilia.FamiliaID != 0 && updatedProductoFamilia.FamiliaID != productoFamilia.FamiliaID {
			// Verificar que la nueva familia existe y pertenece a la misma empresa
			var familia Familia
			if err := tx.Where("id = ? AND empresa_id = ?", updatedProductoFamilia.FamiliaID, empresaID).First(&familia).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("familia con id %d no encontrada o no pertenece a la empresa", updatedProductoFamilia.FamiliaID)
				}
				return fmt.Errorf("error al verificar familia: %w", err)
			}
		}

		if updatedProductoFamilia.ProductoID != 0 && updatedProductoFamilia.ProductoID != productoFamilia.ProductoID {
			// Verificar que el nuevo producto existe y pertenece a la misma empresa
			var producto Producto
			if err := tx.Where("id = ? AND empresa_id = ?", updatedProductoFamilia.ProductoID, empresaID).First(&producto).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("producto con id %d no encontrado o no pertenece a la empresa", updatedProductoFamilia.ProductoID)
				}
				return fmt.Errorf("error al verificar producto: %w", err)
			} // Verificar que no existe ya una relación entre este producto y familia
			var count int64
			var familiaID uint64
			if updatedProductoFamilia.FamiliaID != 0 {
				familiaID = updatedProductoFamilia.FamiliaID
			} else {
				familiaID = productoFamilia.FamiliaID
			}

			if err := tx.Model(&ProductoFamilia{}).
				Where("producto_id = ? AND familia_id = ? AND id != ?",
					updatedProductoFamilia.ProductoID,
					familiaID,
					id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar duplicidad: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("ya existe una relación entre el producto %d y la familia", updatedProductoFamilia.ProductoID)
			}
		}

		return tx.Model(&productoFamilia).Updates(updatedProductoFamilia).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoFamilia, fmt.Errorf("relación producto-familia con id %d no encontrada", id)
		}
		return productoFamilia, fmt.Errorf("error al actualizar relación producto-familia con id %d: %w", id, err)
	}
	return productoFamilia, nil
}

// DeleteProductoFamilia elimina una relación ProductoFamilia
func (r *productoFamiliaRepo) DeleteProductoFamilia(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.
			Joins("JOIN familias ON producto_familias.familia_id = familias.id").
			Where("producto_familias.id = ? AND familias.empresa_id = ?", id, empresaID).
			Delete(&ProductoFamilia{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("relación producto-familia con id %d no encontrada o no pertenece a la empresa", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar relación producto-familia: %w", err)
	}
	return nil
}

// SeedProductoFamilias inicializa múltiples relaciones ProductoFamilia
func (r *productoFamiliaRepo) SeedProductoFamilias(productoFamilias []ProductoFamilia) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&productoFamilias).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar relaciones producto-familia: %w", err)
	}
	return nil
}

// PatchProductoFamilia actualiza parcialmente una relación ProductoFamilia
func (r *productoFamiliaRepo) PatchProductoFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoFamilia, error) {
	var productoFamilia ProductoFamilia
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Primero recuperamos la relación actual para verificar que pertenece a la empresa
		if err := tx.
			Joins("JOIN familias ON producto_familias.familia_id = familias.id").
			Where("producto_familias.id = ? AND familias.empresa_id = ?", id, empresaID).
			First(&productoFamilia).Error; err != nil {
			return err
		}

		// Si se está cambiando la familia, verificar que existe y pertenece a la misma empresa
		if familiaID, ok := fields["familia_id"].(float64); ok {
			var familia Familia
			if err := tx.Where("id = ? AND empresa_id = ?", uint64(familiaID), empresaID).First(&familia).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("familia con id %d no encontrada o no pertenece a la empresa", uint64(familiaID))
				}
				return fmt.Errorf("error al verificar familia: %w", err)
			}
		}

		// Si se está cambiando el producto, verificar que existe y pertenece a la misma empresa
		if productoID, ok := fields["producto_id"].(float64); ok {
			var producto Producto
			if err := tx.Where("id = ? AND empresa_id = ?", uint64(productoID), empresaID).First(&producto).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("producto con id %d no encontrado o no pertenece a la empresa", uint64(productoID))
				}
				return fmt.Errorf("error al verificar producto: %w", err)
			}

			// Verificar que no existe ya una relación entre este producto y familia
			newFamiliaID := productoFamilia.FamiliaID
			if familiaID, ok := fields["familia_id"].(float64); ok {
				newFamiliaID = uint64(familiaID)
			}

			var count int64
			if err := tx.Model(&ProductoFamilia{}).
				Where("producto_id = ? AND familia_id = ? AND id != ?", uint64(productoID), newFamiliaID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("error al verificar duplicidad: %w", err)
			}
			if count > 0 {
				return fmt.Errorf("ya existe una relación entre el producto %d y la familia %d", uint64(productoID), newFamiliaID)
			}
		}

		return tx.Model(&productoFamilia).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return productoFamilia, fmt.Errorf("relación producto-familia con id %d no encontrada", id)
		}
		return productoFamilia, fmt.Errorf("error al actualizar parcialmente relación producto-familia con id %d: %w", id, err)
	}
	return productoFamilia, nil
}
