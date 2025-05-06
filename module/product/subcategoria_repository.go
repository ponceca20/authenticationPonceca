package product

import (
	"errors"
	"fmt"
	"practicev2/database"

	"gorm.io/gorm"
)

// SubCategoriaRepository define la interfaz para operaciones de persistencia de SubCategoria
type SubCategoriaRepository interface {
	GetAllSubCategorias(empresaID uint64) ([]SubCategoria, error)
	GetSubCategoria(id uint64, empresaID uint64) (SubCategoria, error)
	CreateSubCategoria(sub *SubCategoria) error
	UpdateSubCategoria(id uint64, empresaID uint64, sub *SubCategoria) (SubCategoria, error)
	DeleteSubCategoria(id uint64, empresaID uint64) error
	PatchSubCategoria(id uint64, fields map[string]interface{}) (SubCategoria, error)
}

// subCategoriaRepo es la implementación con GORM
type subCategoriaRepo struct {
	DB *gorm.DB
}

// NewSubCategoriaRepository crea una instancia del repositorio
func NewSubCategoriaRepository() SubCategoriaRepository {
	return &subCategoriaRepo{DB: database.DBconn}
}

// GetAllSubCategorias lista todas las subcategorias de las categorias de una empresa
func (r *subCategoriaRepo) GetAllSubCategorias(empresaID uint64) ([]SubCategoria, error) {
	var subs []SubCategoria
	err := r.DB.Joins("JOIN categoria ON categoria.id = sub_categoria.categoria_id AND categoria.empresa_id = ?", empresaID).
		Preload("Categoria").Find(&subs).Error
	if err != nil {
		return nil, fmt.Errorf("error al recuperar subcategorias: %w", err)
	}
	return subs, nil
}

// GetSubCategoria obtiene una subcategoria por ID y empresa mediante join
func (r *subCategoriaRepo) GetSubCategoria(id uint64, empresaID uint64) (SubCategoria, error) {
	var sub SubCategoria
	err := r.DB.Joins("JOIN categoria ON categoria.id = sub_categoria.categoria_id AND categoria.empresa_id = ?", empresaID).
		Preload("Categoria").Where("sub_categoria.id = ?", id).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sub, fmt.Errorf("subcategoria con id %d no encontrada", id)
		}
		return sub, fmt.Errorf("error al recuperar subcategoria con id %d: %w", id, err)
	}
	return sub, nil
}

// CreateSubCategoria crea una nueva subcategoria
func (r *subCategoriaRepo) CreateSubCategoria(sub *SubCategoria) error {
	if err := r.DB.Create(sub).Error; err != nil {
		return fmt.Errorf("error al crear subcategoria: %w", err)
	}
	return nil
}

// UpdateSubCategoria actualiza una subcategoria existente
func (r *subCategoriaRepo) UpdateSubCategoria(id uint64, empresaID uint64, updated *SubCategoria) (SubCategoria, error) {
	var sub SubCategoria
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN categoria ON categoria.id = sub_categoria.categoria_id AND categoria.empresa_id = ?", empresaID).
			Where("sub_categoria.id = ?", id).First(&sub).Error; err != nil {
			return err
		}
		return tx.Model(&sub).Updates(updated).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sub, fmt.Errorf("subcategoria con id %d no encontrada", id)
		}
		return sub, fmt.Errorf("error al actualizar subcategoria con id %d: %w", id, err)
	}
	return sub, nil
}

// DeleteSubCategoria elimina una subcategoria, previene si hay productos asociados
func (r *subCategoriaRepo) DeleteSubCategoria(id uint64, empresaID uint64) error {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Verificar productos asociados
		var count int64
		if err := tx.Model(&Producto{}).Where("sub_categoria_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("no se puede eliminar subcategoria porque tiene %d productos asociados", count)
		}
		res := tx.Joins("JOIN categoria ON categoria.id = sub_categoria.categoria_id AND categoria.empresa_id = ?", empresaID).
			Where("sub_categoria.id = ?", id).Delete(&SubCategoria{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("subcategoria con id %d no encontrada", id)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error al eliminar subcategoria: %w", err)
	}
	return nil
}

// PatchSubCategoria actualiza parcialmente una subcategoria
func (r *subCategoriaRepo) PatchSubCategoria(id uint64, fields map[string]interface{}) (SubCategoria, error) {
	var sub SubCategoria
	delete(fields, "empresa_id")
	// El parámetro empresaID se elimina, toda la lógica de empresa se maneja por el JOIN con categoria
	
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Joins("JOIN categoria ON categoria.id = sub_categoria.categoria_id").
			Where("sub_categoria.id = ?", id).First(&sub).Error; err != nil {
			return err
		}
		return tx.Model(&sub).Updates(fields).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sub, fmt.Errorf("subcategoria con id %d no encontrada", id)
		}
		return sub, fmt.Errorf("error al actualizar parcialmente subcategoria con id %d: %w", id, err)
	}
	return sub, nil
}
