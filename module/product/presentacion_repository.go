package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PresentacionRepository define la interfaz para operaciones de persistencia de Presentaciones
type PresentacionRepository interface {
	GetAllPresentaciones(productoID uint64) ([]Presentacion, error)
	GetPresentacion(id uint64) (Presentacion, error)
	CreatePresentacion(presentacion *Presentacion) error
	UpdatePresentacion(id uint64, presentacion *Presentacion) (Presentacion, error)
	DeletePresentacion(id uint64) error
	SeedPresentaciones(presentaciones []Presentacion) error
	PatchPresentacion(id uint64, fields map[string]interface{}) (Presentacion, error)
	GetAllPresentacionesByEmpresa(empresaID uint64) ([]Presentacion, error)
}

type presentacionRepo struct {
	DB *gorm.DB
}

// NewPresentacionRepository crea una instancia del repositorio
func NewPresentacionRepository() PresentacionRepository {
	return &presentacionRepo{DB: database.DBconn}
}

// GetAllPresentaciones obtiene todas las presentaciones para un producto
func (r *presentacionRepo) GetAllPresentaciones(productoID uint64) ([]Presentacion, error) {
	var presentaciones []Presentacion
	err := r.DB.Where("producto_id = ?", productoID).Find(&presentaciones).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar presentaciones: %w", err)
	}
	return presentaciones, nil
}

// GetPresentacion obtiene una presentación específica
func (r *presentacionRepo) GetPresentacion(id uint64) (Presentacion, error) {
	var presentacion Presentacion
	err := r.DB.Where("id = ?", id).First(&presentacion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presentacion, fmt.Errorf("presentación con id %d no encontrada", id)
		}
		return presentacion, fmt.Errorf("error al recuperar presentación con id %d: %w", id, err)
	}
	return presentacion, nil
}

// CreatePresentacion crea una nueva presentación
func (r *presentacionRepo) CreatePresentacion(presentacion *Presentacion) error {
	return r.DB.Create(presentacion).Error
}

// UpdatePresentacion actualiza una presentación existente
func (r *presentacionRepo) UpdatePresentacion(id uint64, presentacion *Presentacion) (Presentacion, error) {
	var existing Presentacion
	err := r.DB.Where("id = ?", id).First(&existing).Error
	if err != nil {
		return existing, fmt.Errorf("presentación con id %d no encontrada: %w", id, err)
	}
	if err := r.DB.Model(&existing).Updates(presentacion).Error; err != nil {
		return existing, fmt.Errorf("error al actualizar presentación con id %d: %w", id, err)
	}
	return existing, nil
}

// DeletePresentacion elimina una presentación
func (r *presentacionRepo) DeletePresentacion(id uint64) error {
	result := r.DB.Where("id = ?", id).Delete(&Presentacion{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("presentación con id %d no encontrada", id)
	}
	return nil
}

// SeedPresentaciones inicializa múltiples presentaciones
func (r *presentacionRepo) SeedPresentaciones(presentaciones []Presentacion) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&presentaciones).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar presentaciones: %w", err)
	}
	return nil
}

// PatchPresentacion actualiza parcialmente una presentación
func (r *presentacionRepo) PatchPresentacion(id uint64, fields map[string]interface{}) (Presentacion, error) {
	var presentacion Presentacion
	delete(fields, "empresa_id")
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&presentacion).Error; err != nil {
			return err
		}
		// Si se cambia el producto_id, verificar que el nuevo producto exista
		if newProductoID, ok := fields["producto_id"].(float64); ok {
			var count int64
			if err := tx.Model(&Producto{}).
				Where("id = ?", uint64(newProductoID)).
				Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("el producto nuevo con ID %d no existe", uint64(newProductoID))
			}
		}
		return tx.Model(&presentacion).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presentacion, fmt.Errorf("presentación con id %d no encontrada", id)
		}
		return presentacion, fmt.Errorf("error al actualizar parcialmente presentación con id %d: %w", id, err)
	}
	return presentacion, nil
}

// GetAllPresentacionesByEmpresa obtiene todas las presentaciones filtradas por empresa_id
func (r *presentacionRepo) GetAllPresentacionesByEmpresa(empresaID uint64) ([]Presentacion, error) {
	var presentaciones []Presentacion
	err := r.DB.Joins("JOIN producto ON producto.id = presentacion.producto_id").Where("producto.empresa_id = ?", empresaID).Find(&presentaciones).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar presentaciones por empresa: %w", err)
	}
	return presentaciones, nil
}
