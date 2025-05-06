package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// PresentacionCompuestaRepository define la interfaz para operar con PresentacionCompuesta
type PresentacionCompuestaRepository interface {
	GetAllPresentacionCompuestas(empresaID uint64) ([]PresentacionCompuesta, error)
	GetPresentacionCompuesta(id uint64, empresaID uint64) (PresentacionCompuesta, error)
	CreatePresentacionCompuesta(pc *PresentacionCompuesta) error
	UpdatePresentacionCompuesta(id uint64, empresaID uint64, pc *PresentacionCompuesta) (PresentacionCompuesta, error)
	DeletePresentacionCompuesta(id uint64, empresaID uint64) error
	SeedPresentacionCompuestas(pcs []PresentacionCompuesta) error
	PatchPresentacionCompuesta(id uint64, empresaID uint64, fields map[string]interface{}) (PresentacionCompuesta, error)
}

type presentacionCompuestaRepo struct {
	DB *gorm.DB
}

// NewPresentacionCompuestaRepository crea una instancia del repositorio
func NewPresentacionCompuestaRepository() PresentacionCompuestaRepository {
	return &presentacionCompuestaRepo{DB: database.DBconn}
}

func (r *presentacionCompuestaRepo) GetAllPresentacionCompuestas(empresaID uint64) ([]PresentacionCompuesta, error) {
	var pcs []PresentacionCompuesta
	// Se une Presentacion y Producto para validar la multiempresa vía producto.empresa_id
	err := r.DB.Joins("JOIN presentacion ON presentacion.id = presentacion_compuesta.presentacion_principal_id").
		Joins("JOIN producto ON producto.id = presentacion.producto_id").
		Where("producto.empresa_id = ?", empresaID).
		Find(&pcs).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar presentacion compuesta: %w", err)
	}
	return pcs, nil
}

func (r *presentacionCompuestaRepo) GetPresentacionCompuesta(id uint64, empresaID uint64) (PresentacionCompuesta, error) {
	var pc PresentacionCompuesta
	err := r.DB.Joins("JOIN presentacion ON presentacion.id = presentacion_compuesta.presentacion_principal_id").
		Joins("JOIN producto ON producto.id = presentacion.producto_id").
		Where("presentacion_compuesta.id = ? AND producto.empresa_id = ?", id, empresaID).
		First(&pc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pc, fmt.Errorf("presentacion compuesta con id %d no encontrada", id)
		}
		return pc, fmt.Errorf("error al recuperar presentacion compuesta con id %d: %w", id, err)
	}
	return pc, nil
}

func (r *presentacionCompuestaRepo) CreatePresentacionCompuesta(pc *PresentacionCompuesta) error {
	err := r.DB.Create(pc).Error
	if err != nil {
		return fmt.Errorf("error al crear presentacion compuesta: %w", err)
	}
	return nil
}

func (r *presentacionCompuestaRepo) UpdatePresentacionCompuesta(id uint64, empresaID uint64, updatedPC *PresentacionCompuesta) (PresentacionCompuesta, error) {
	var pc PresentacionCompuesta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN presentacion ON presentacion.id = presentacion_compuesta.presentacion_principal_id").
			Joins("JOIN producto ON producto.id = presentacion.producto_id").
			Where("presentacion_compuesta.id = ? AND producto.empresa_id = ?", id, empresaID).
			First(&pc).Error; err != nil {
			return err
		}
		return tx.Model(&pc).Updates(updatedPC).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pc, fmt.Errorf("presentacion compuesta con id %d no encontrada", id)
		}
		return pc, fmt.Errorf("error al actualizar presentacion compuesta con id %d: %w", id, err)
	}
	return pc, nil
}

func (r *presentacionCompuestaRepo) DeletePresentacionCompuesta(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Joins("JOIN presentacion ON presentacion.id = presentacion_compuesta.presentacion_principal_id").
			Joins("JOIN producto ON producto.id = presentacion.producto_id").
			Where("presentacion_compuesta.id = ? AND producto.empresa_id = ?", id, empresaID).
			Delete(&PresentacionCompuesta{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("presentacion compuesta con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar presentacion compuesta: %w", err)
	}
	return nil
}

func (r *presentacionCompuestaRepo) SeedPresentacionCompuestas(pcs []PresentacionCompuesta) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&pcs).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar presentacion compuesta: %w", err)
	}
	return nil
}

func (r *presentacionCompuestaRepo) PatchPresentacionCompuesta(id uint64, empresaID uint64, fields map[string]interface{}) (PresentacionCompuesta, error) {
	var pc PresentacionCompuesta
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN presentacion ON presentacion.id = presentacion_compuesta.presentacion_principal_id").
			Joins("JOIN producto ON producto.id = presentacion.producto_id").
			Where("presentacion_compuesta.id = ? AND producto.empresa_id = ?", id, empresaID).
			First(&pc).Error; err != nil {
			return err
		}
		return tx.Model(&pc).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pc, fmt.Errorf("presentacion compuesta con id %d no encontrada", id)
		}
		return pc, fmt.Errorf("error al actualizar parcialmente presentacion compuesta con id %d: %w", id, err)
	}
	return pc, nil
}
