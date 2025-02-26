package users

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// UsuarioRepository define la interfaz para operaciones de persistencia
type UsuarioRepository interface {
	GetAllUsuarios() ([]Usuario, error)
	GetUsuario(id uint64) (Usuario, error)
	CreateUsuario(usuario *Usuario) error
	UpdateUsuario(id uint64, updatedUsuario *Usuario) (Usuario, error)
	DeleteUsuario(id uint64) error
	SeedUsuarios(usuarios []Usuario) error
	PatchUsuario(id uint64, fields map[string]interface{}) (Usuario, error)
}

type usuarioRepo struct {
	DB *gorm.DB
}

// NewUsuarioRepository crea una nueva instancia del repositorio
func NewUsuarioRepository() UsuarioRepository {
	return &usuarioRepo{DB: database.DBconn}
}

func (r *usuarioRepo) GetAllUsuarios() ([]Usuario, error) {
	var usuarios []Usuario
	err := r.DB.Find(&usuarios).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar usuarios: %w", err)
	}
	return usuarios, nil
}

func (r *usuarioRepo) GetUsuario(id uint64) (Usuario, error) {
	var usuario Usuario
	err := r.DB.First(&usuario, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usuario, fmt.Errorf("usuario con id %d no encontrado", id)
		}
		return usuario, fmt.Errorf("error al recuperar usuario con id %d: %w", id, err)
	}
	return usuario, nil
}

func (r *usuarioRepo) CreateUsuario(usuario *Usuario) error {
	err := r.DB.Create(usuario).Error
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}
	return nil
}

func (r *usuarioRepo) UpdateUsuario(id uint64, updatedUsuario *Usuario) (Usuario, error) {
	var usuario Usuario
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&usuario, id).Error; err != nil {
			return err
		}
		return tx.Model(&usuario).Updates(updatedUsuario).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usuario, fmt.Errorf("usuario con id %d no encontrado", id)
		}
		return usuario, fmt.Errorf("error al actualizar usuario con id %d: %w", id, err)
	}
	return usuario, nil
}

func (r *usuarioRepo) DeleteUsuario(id uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Usuario{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("usuario con id %d no encontrado", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar usuario: %w", err)
	}
	return nil
}

func (r *usuarioRepo) SeedUsuarios(usuarios []Usuario) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&usuarios).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar usuarios: %w", err)
	}
	return nil
}

func (r *usuarioRepo) PatchUsuario(id uint64, fields map[string]interface{}) (Usuario, error) {
	var usuario Usuario
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&usuario, id).Error; err != nil {
			return err
		}
		return tx.Model(&usuario).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usuario, fmt.Errorf("usuario con id %d no encontrado", id)
		}
		return usuario, fmt.Errorf("error al actualizar parcialmente usuario con id %d: %w", id, err)
	}
	return usuario, nil
}
