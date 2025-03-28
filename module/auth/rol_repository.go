package auth

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// RolRepository define la interfaz para operaciones en Rol.
type RolRepository interface {
	GetAllRoles() ([]Rol, error)
	GetRol(id uint64) (Rol, error)
	CreateRol(rol *Rol) error
	UpdateRol(id uint64, updatedRol *Rol) (Rol, error)
	DeleteRol(id uint64) error
	SeedRoles(roles []Rol) error
	PatchRol(id uint64, fields map[string]interface{}) (Rol, error)
}

type rolRepo struct {
	DB *gorm.DB
}

func NewRolRepository() RolRepository {
	return &rolRepo{DB: database.DBconn}
}

func (r *rolRepo) GetAllRoles() ([]Rol, error) {
	var roles []Rol
	err := r.DB.Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar roles: %w", err)
	}
	return roles, nil
}

func (r *rolRepo) GetRol(id uint64) (Rol, error) {
	var rol Rol
	err := r.DB.First(&rol, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rol, fmt.Errorf("rol con id %d no encontrado", id)
		}
		return rol, fmt.Errorf("error al recuperar rol con id %d: %w", id, err)
	}
	return rol, nil
}

func (r *rolRepo) CreateRol(rol *Rol) error {
	err := r.DB.Create(rol).Error
	if err != nil {
		return fmt.Errorf("error al crear rol: %w", err)
	}
	return nil
}

func (r *rolRepo) UpdateRol(id uint64, updatedRol *Rol) (Rol, error) {
	var rol Rol
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&rol, id).Error; err != nil {
			return err
		}
		return tx.Model(&rol).Updates(updatedRol).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rol, fmt.Errorf("rol con id %d no encontrado", id)
		}
		return rol, fmt.Errorf("error al actualizar rol con id %d: %w", id, err)
	}
	return rol, nil
}

func (r *rolRepo) DeleteRol(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Rol{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("rol con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar rol: %w", err)
	}
	return nil
}

func (r *rolRepo) SeedRoles(roles []Rol) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&roles).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar roles: %w", err)
	}
	return nil
}

func (r *rolRepo) PatchRol(id uint64, fields map[string]interface{}) (Rol, error) {
	var rol Rol
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&rol, id).Error; err != nil {
			return err
		}
		return tx.Model(&rol).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rol, fmt.Errorf("rol con id %d no encontrado", id)
		}
		return rol, fmt.Errorf("error al actualizar parcialmente rol con id %d: %w", id, err)
	}
	return rol, nil
}
