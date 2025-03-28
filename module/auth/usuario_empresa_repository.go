package auth

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// UsuarioEmpresaRepository define la interfaz para operaciones en UsuarioEmpresa.
type UsuarioEmpresaRepository interface {
	GetAllUsuarioEmpresas() ([]UsuarioEmpresa, error)
	GetUsuarioEmpresa(id uint64) (UsuarioEmpresa, error)
	CreateUsuarioEmpresa(ue *UsuarioEmpresa) error
	UpdateUsuarioEmpresa(id uint64, updatedUE *UsuarioEmpresa) (UsuarioEmpresa, error)
	DeleteUsuarioEmpresa(id uint64) error
	SeedUsuarioEmpresas(ues []UsuarioEmpresa) error
	PatchUsuarioEmpresa(id uint64, fields map[string]interface{}) (UsuarioEmpresa, error)
}

type usuarioEmpresaRepo struct {
	DB *gorm.DB
}

func NewUsuarioEmpresaRepository() UsuarioEmpresaRepository {
	return &usuarioEmpresaRepo{DB: database.DBconn}
}

func (r *usuarioEmpresaRepo) GetAllUsuarioEmpresas() ([]UsuarioEmpresa, error) {
	var ues []UsuarioEmpresa
	err := r.DB.Find(&ues).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar usuario empresas: %w", err)
	}
	return ues, nil
}

func (r *usuarioEmpresaRepo) GetUsuarioEmpresa(id uint64) (UsuarioEmpresa, error) {
	var ue UsuarioEmpresa
	err := r.DB.First(&ue, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ue, fmt.Errorf("usuario empresa con id %d no encontrada", id)
		}
		return ue, fmt.Errorf("error al recuperar usuario empresa con id %d: %w", id, err)
	}
	return ue, nil
}

func (r *usuarioEmpresaRepo) CreateUsuarioEmpresa(ue *UsuarioEmpresa) error {
	err := r.DB.Create(ue).Error
	if err != nil {
		return fmt.Errorf("error al crear usuario empresa: %w", err)
	}
	return nil
}

func (r *usuarioEmpresaRepo) UpdateUsuarioEmpresa(id uint64, updatedUE *UsuarioEmpresa) (UsuarioEmpresa, error) {
	var ue UsuarioEmpresa
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&ue, id).Error; err != nil {
			return err
		}
		return tx.Model(&ue).Updates(updatedUE).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ue, fmt.Errorf("usuario empresa con id %d no encontrada", id)
		}
		return ue, fmt.Errorf("error al actualizar usuario empresa con id %d: %w", id, err)
	}
	return ue, nil
}

func (r *usuarioEmpresaRepo) DeleteUsuarioEmpresa(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&UsuarioEmpresa{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("usuario empresa con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar usuario empresa: %w", err)
	}
	return nil
}

func (r *usuarioEmpresaRepo) SeedUsuarioEmpresas(ues []UsuarioEmpresa) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ues).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar usuario empresas: %w", err)
	}
	return nil
}

func (r *usuarioEmpresaRepo) PatchUsuarioEmpresa(id uint64, fields map[string]interface{}) (UsuarioEmpresa, error) {
	var ue UsuarioEmpresa
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&ue, id).Error; err != nil {
			return err
		}
		return tx.Model(&ue).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ue, fmt.Errorf("usuario empresa con id %d no encontrada", id)
		}
		return ue, fmt.Errorf("error al actualizar parcialmente usuario empresa con id %d: %w", id, err)
	}
	return ue, nil
}
