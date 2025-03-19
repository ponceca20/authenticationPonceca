package users

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// RolModuloRepository define la interfaz para operaciones de persistencia de RolModulo
type RolModuloRepository interface {
	GetAllRolModulos() ([]RolModulo, error)
	GetRolModulo(id uint64) (RolModulo, error)
	CreateRolModulo(rolmodulo *RolModulo) error
	UpdateRolModulo(id uint64, updatedRolModulo *RolModulo) (RolModulo, error)
	DeleteRolModulo(id uint64) error
	SeedRolModulos(rolmodulos []RolModulo) error
	PatchRolModulo(id uint64, fields map[string]interface{}) (RolModulo, error)
}

type rolModuloRepo struct {
	DB *gorm.DB
}

// NewRolModuloRepository crea una nueva instancia del repositorio para RolModulo
func NewRolModuloRepository() RolModuloRepository {
	return &rolModuloRepo{DB: database.DBconn}
}

func (r *rolModuloRepo) GetAllRolModulos() ([]RolModulo, error) {
	var rolmodulos []RolModulo
	err := r.DB.Find(&rolmodulos).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar rolmodulos: %w", err)
	}
	return rolmodulos, nil
}

func (r *rolModuloRepo) GetRolModulo(id uint64) (RolModulo, error) {
	var rolmodulo RolModulo
	err := r.DB.First(&rolmodulo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rolmodulo, fmt.Errorf("rolmodulo con id %d no encontrado", id)
		}
		return rolmodulo, fmt.Errorf("error al recuperar rolmodulo con id %d: %w", id, err)
	}
	return rolmodulo, nil
}

func (r *rolModuloRepo) CreateRolModulo(rolmodulo *RolModulo) error {
	err := r.DB.Create(rolmodulo).Error
	if err != nil {
		return fmt.Errorf("error al crear rolmodulo: %w", err)
	}
	return nil
}

func (r *rolModuloRepo) UpdateRolModulo(id uint64, updatedRolModulo *RolModulo) (RolModulo, error) {
	var rolmodulo RolModulo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&rolmodulo, id).Error; err != nil {
			return err
		}
		return tx.Model(&rolmodulo).Updates(updatedRolModulo).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rolmodulo, fmt.Errorf("rolmodulo con id %d no encontrado", id)
		}
		return rolmodulo, fmt.Errorf("error al actualizar rolmodulo con id %d: %w", id, err)
	}
	return rolmodulo, nil
}

func (r *rolModuloRepo) DeleteRolModulo(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&RolModulo{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("rolmodulo con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar rolmodulo: %w", err)
	}
	return nil
}

func (r *rolModuloRepo) SeedRolModulos(rolmodulos []RolModulo) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&rolmodulos).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar rolmodulos: %w", err)
	}
	return nil
}

func (r *rolModuloRepo) PatchRolModulo(id uint64, fields map[string]interface{}) (RolModulo, error) {
	var rolmodulo RolModulo
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&rolmodulo, id).Error; err != nil {
			return err
		}
		return tx.Model(&rolmodulo).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rolmodulo, fmt.Errorf("rolmodulo con id %d no encontrado", id)
		}
		return rolmodulo, fmt.Errorf("error al actualizar parcialmente rolmodulo con id %d: %w", id, err)
	}
	return rolmodulo, nil
}
