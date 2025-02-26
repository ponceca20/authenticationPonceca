package users

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// SesionRepository define la interfaz para operaciones de persistencia en Sesion
type SesionRepository interface {
	GetAllSesions() ([]Sesion, error)
	GetSesion(id uint64) (Sesion, error)
	CreateSesion(sesion *Sesion) error
	UpdateSesion(id uint64, updatedSesion *Sesion) (Sesion, error)
	DeleteSesion(id uint64) error
	SeedSesions(sesiones []Sesion) error
	PatchSesion(id uint64, fields map[string]interface{}) (Sesion, error)
}

type sesionRepo struct {
	DB *gorm.DB
}

// NewSesionRepository crea una nueva instancia del repositorio
func NewSesionRepository() SesionRepository {
	return &sesionRepo{DB: database.DBconn}
}

func (r *sesionRepo) GetAllSesions() ([]Sesion, error) {
	var sesions []Sesion
	err := r.DB.Find(&sesions).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar sesiones: %w", err)
	}
	return sesions, nil
}

func (r *sesionRepo) GetSesion(id uint64) (Sesion, error) {
	var sesion Sesion
	err := r.DB.First(&sesion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sesion, fmt.Errorf("sesión con id %d no encontrada", id)
		}
		return sesion, fmt.Errorf("error al recuperar sesión con id %d: %w", id, err)
	}
	return sesion, nil
}

func (r *sesionRepo) CreateSesion(sesion *Sesion) error {
	err := r.DB.Create(sesion).Error
	if err != nil {
		return fmt.Errorf("error al crear sesión: %w", err)
	}
	return nil
}

func (r *sesionRepo) UpdateSesion(id uint64, updatedSesion *Sesion) (Sesion, error) {
	var sesion Sesion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&sesion, id).Error; err != nil {
			return err
		}
		return tx.Model(&sesion).Updates(updatedSesion).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sesion, fmt.Errorf("sesión con id %d no encontrada", id)
		}
		return sesion, fmt.Errorf("error al actualizar sesión con id %d: %w", id, err)
	}
	return sesion, nil
}

func (r *sesionRepo) DeleteSesion(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Sesion{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("sesión con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar sesión: %w", err)
	}
	return nil
}

func (r *sesionRepo) SeedSesions(sesiones []Sesion) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&sesiones).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar sesiones: %w", err)
	}
	return nil
}

func (r *sesionRepo) PatchSesion(id uint64, fields map[string]interface{}) (Sesion, error) {
	var sesion Sesion
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&sesion, id).Error; err != nil {
			return err
		}
		return tx.Model(&sesion).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sesion, fmt.Errorf("sesión con id %d no encontrada", id)
		}
		return sesion, fmt.Errorf("error al actualizar parcialmente sesión con id %d: %w", id, err)
	}
	return sesion, nil
}
