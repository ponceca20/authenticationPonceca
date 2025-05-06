package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ClaseRepository define operaciones de persistencia para Clase
type ClaseRepository interface {
	GetAllClases(empresaID uint64) ([]Clase, error)
	GetClase(id uint64, empresaID uint64) (Clase, error)
	CreateClase(clase *Clase) error
	UpdateClase(id uint64, empresaID uint64, clase *Clase) (Clase, error)
	DeleteClase(id uint64, empresaID uint64) error
	PatchClase(id uint64, empresaID uint64, fields map[string]interface{}) (Clase, error)
}

type claseRepo struct {
	DB *gorm.DB
}

// NewClaseRepository crea una instancia del repositorio
func NewClaseRepository() ClaseRepository {
	return &claseRepo{DB: database.DBconn}
}

// GetAllClases lista todas las clases de una empresa
func (r *claseRepo) GetAllClases(empresaID uint64) ([]Clase, error) {
	var clases []Clase
	err := r.DB.Where("empresa_id = ?", empresaID).Find(&clases).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar clases: %w", err)
	}
	return clases, nil
}

// GetClase obtiene una clase por ID y empresa
func (r *claseRepo) GetClase(id uint64, empresaID uint64) (Clase, error) {
	var clase Clase
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).First(&clase).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return clase, fmt.Errorf("clase con id %d no encontrada", id)
		}
		return clase, fmt.Errorf("error al recuperar clase con id %d: %w", id, err)
	}
	return clase, nil
}

// CreateClase crea una nueva clase
func (r *claseRepo) CreateClase(clase *Clase) error {
	if err := r.DB.Create(clase).Error; err != nil {
		return fmt.Errorf("error al crear clase: %w", err)
	}
	return nil
}

// UpdateClase actualiza una clase existente
func (r *claseRepo) UpdateClase(id uint64, empresaID uint64, updated *Clase) (Clase, error) {
	var clase Clase
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&clase).Error; err != nil {
			return err
		}
		return tx.Model(&clase).Updates(updated).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return clase, fmt.Errorf("clase con id %d no encontrada", id)
		}
		return clase, fmt.Errorf("error al actualizar clase con id %d: %w", id, err)
	}
	return clase, nil
}

// DeleteClase elimina una clase
func (r *claseRepo) DeleteClase(id uint64, empresaID uint64) error {
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Clase{}).Error
	if err != nil {
		return fmt.Errorf("error al eliminar clase: %w", err)
	}
	return nil
}

// PatchClase actualiza parcialmente una clase
func (r *claseRepo) PatchClase(id uint64, empresaID uint64, fields map[string]interface{}) (Clase, error) {
	var clase Clase
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&clase).Error; err != nil {
			return err
		}
		return tx.Model(&clase).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return clase, fmt.Errorf("clase con id %d no encontrada", id)
		}
		return clase, fmt.Errorf("error al actualizar parcialmente clase con id %d: %w", id, err)
	}
	return clase, nil
}
