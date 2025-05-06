package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// CategoriaRepository define las operaciones de persistencia para Categoria
type CategoriaRepository interface {
	GetAllCategorias(empresaID uint64) ([]Categoria, error)
	GetCategoria(id uint64, empresaID uint64) (Categoria, error)
	CreateCategoria(categoria *Categoria) error
	UpdateCategoria(id uint64, empresaID uint64, categoria *Categoria) (Categoria, error)
	DeleteCategoria(id uint64, empresaID uint64) error
	PatchCategoria(id uint64, empresaID uint64, fields map[string]interface{}) (Categoria, error)
}

type categoriaRepo struct {
	DB *gorm.DB
}

// NewCategoriaRepository crea una instancia del repositorio
func NewCategoriaRepository() CategoriaRepository {
	return &categoriaRepo{DB: database.DBconn}
}

// GetAllCategorias lista todas las categorias de una empresa
func (r *categoriaRepo) GetAllCategorias(empresaID uint64) ([]Categoria, error) {
	var cats []Categoria
	err := r.DB.Where("empresa_id = ?", empresaID).Preload("SubCategorias").Find(&cats).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar categorias: %w", err)
	}
	return cats, nil
}

// GetCategoria obtiene una categoria por ID y empresa
func (r *categoriaRepo) GetCategoria(id uint64, empresaID uint64) (Categoria, error) {
	var cat Categoria
	err := r.DB.Where("id = ? AND empresa_id = ?", id, empresaID).Preload("SubCategorias").First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cat, fmt.Errorf("categoria con id %d no encontrada", id)
		}
		return cat, fmt.Errorf("error al recuperar categoria con id %d: %w", id, err)
	}
	return cat, nil
}

// CreateCategoria crea una nueva categoria
func (r *categoriaRepo) CreateCategoria(categoria *Categoria) error {
	if err := r.DB.Create(categoria).Error; err != nil {
		return fmt.Errorf("error al crear categoria: %w", err)
	}
	return nil
}

// UpdateCategoria actualiza una categoria existente
func (r *categoriaRepo) UpdateCategoria(id uint64, empresaID uint64, updated *Categoria) (Categoria, error) {
	var cat Categoria
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&cat).Error; err != nil {
			return err
		}
		return tx.Model(&cat).Updates(updated).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cat, fmt.Errorf("categoria con id %d no encontrada", id)
		}
		return cat, fmt.Errorf("error al actualizar categoria con id %d: %w", id, err)
	}
	return cat, nil
}

// DeleteCategoria elimina una categoria
func (r *categoriaRepo) DeleteCategoria(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar subcategorias asociadas
		var count int64
		if err := tx.Model(&SubCategoria{}).Where("categoria_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("no se puede eliminar categoria porque tiene %d subcategorias asociadas", count)
		}
		res := tx.Where("id = ? AND empresa_id = ?", id, empresaID).Delete(&Categoria{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("categoria con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar categoria: %w", err)
	}
	return nil
}

// PatchCategoria actualiza parcialmente una categoria
func (r *categoriaRepo) PatchCategoria(id uint64, empresaID uint64, fields map[string]interface{}) (Categoria, error) {
	var cat Categoria
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND empresa_id = ?", id, empresaID).First(&cat).Error; err != nil {
			return err
		}
		return tx.Model(&cat).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cat, fmt.Errorf("categoria con id %d no encontrada", id)
		}
		return cat, fmt.Errorf("error al actualizar parcialmente categoria con id %d: %w", id, err)
	}
	return cat, nil
}
