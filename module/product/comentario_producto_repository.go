package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// ComentarioProductoRepository define la interfaz para operaciones de persistencia de ComentarioProducto
type ComentarioProductoRepository interface {
	GetAllComentariosProducto(productoID *uint64) ([]ComentarioProducto, error)
	GetComentariosGenericos() ([]ComentarioProducto, error)
	GetComentarioProducto(id uint64) (ComentarioProducto, error)
	CreateComentarioProducto(comentario *ComentarioProducto) error
	UpdateComentarioProducto(id uint64, comentario *ComentarioProducto) (ComentarioProducto, error)
	DeleteComentarioProducto(id uint64) error
	SeedComentariosProducto(comentarios []ComentarioProducto) error
	PatchComentarioProducto(id uint64, fields map[string]interface{}) (ComentarioProducto, error)
}

type comentarioProductoRepo struct {
	DB *gorm.DB
}

// NewComentarioProductoRepository crea una instancia del repositorio
func NewComentarioProductoRepository() ComentarioProductoRepository {
	return &comentarioProductoRepo{DB: database.DBconn}
}

// GetAllComentariosProducto obtiene todos los comentarios para un producto específico
// Si productoID es nil, retorna todos los comentarios
func (r *comentarioProductoRepo) GetAllComentariosProducto(productoID *uint64) ([]ComentarioProducto, error) {
	var comentarios []ComentarioProducto
	query := r.DB

	if productoID != nil {
		query = query.Where("producto_id = ?", *productoID)
	}

	err := query.Find(&comentarios).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar comentarios: %w", err)
	}
	return comentarios, nil
}

// GetComentariosGenericos obtiene todos los comentarios genéricos (no asociados a un producto específico)
func (r *comentarioProductoRepo) GetComentariosGenericos() ([]ComentarioProducto, error) {
	var comentarios []ComentarioProducto
	err := r.DB.Where("es_comentario_generico = ?", true).Find(&comentarios).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar comentarios genéricos: %w", err)
	}
	return comentarios, nil
}

// GetComentarioProducto obtiene un comentario específico
func (r *comentarioProductoRepo) GetComentarioProducto(id uint64) (ComentarioProducto, error) {
	var comentario ComentarioProducto
	err := r.DB.Where("id = ?", id).First(&comentario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comentario, fmt.Errorf("comentario con id %d no encontrado", id)
		}
		return comentario, fmt.Errorf("error al recuperar comentario con id %d: %w", id, err)
	}
	return comentario, nil
}

// CreateComentarioProducto crea un nuevo comentario
func (r *comentarioProductoRepo) CreateComentarioProducto(comentario *ComentarioProducto) error {
	err := r.DB.Create(comentario).Error
	if err != nil {
		return fmt.Errorf("error al crear comentario: %w", err)
	}
	return nil
}

// UpdateComentarioProducto actualiza un comentario existente
func (r *comentarioProductoRepo) UpdateComentarioProducto(id uint64, updatedComentario *ComentarioProducto) (ComentarioProducto, error) {
	var comentario ComentarioProducto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&comentario).Error; err != nil {
			return err
		}
		return tx.Model(&comentario).Updates(updatedComentario).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comentario, fmt.Errorf("comentario con id %d no encontrado", id)
		}
		return comentario, fmt.Errorf("error al actualizar comentario con id %d: %w", id, err)
	}
	return comentario, nil
}

// DeleteComentarioProducto elimina un comentario
func (r *comentarioProductoRepo) DeleteComentarioProducto(id uint64) error {
	result := r.DB.Delete(&ComentarioProducto{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar comentario: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("comentario con id %d no encontrado", id)
	}
	return nil
}

// SeedComentariosProducto inicializa múltiples comentarios
func (r *comentarioProductoRepo) SeedComentariosProducto(comentarios []ComentarioProducto) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&comentarios).Error
	})
	if err != nil {
		return fmt.Errorf("error al inicializar comentarios: %w", err)
	}
	return nil
}

// PatchComentarioProducto actualiza parcialmente un comentario
func (r *comentarioProductoRepo) PatchComentarioProducto(id uint64, fields map[string]interface{}) (ComentarioProducto, error) {
	var comentario ComentarioProducto
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&comentario).Error; err != nil {
			return err
		}
		return tx.Model(&comentario).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return comentario, fmt.Errorf("comentario con id %d no encontrado", id)
		}
		return comentario, fmt.Errorf("error al actualizar parcialmente comentario con id %d: %w", id, err)
	}
	return comentario, nil
}
