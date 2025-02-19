package categorias

import (
	"practicev2/database"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAllCategories() ([]Categoria, error)
	GetCategory(id int) (Categoria, error)
	CreateCategory(category *Categoria) error
	UpdateCategory(id int, updatedCategory *Categoria) (Categoria, error)
	DeleteCategory(id int) error
	SeedCategories(categories []Categoria) error
}
type repository struct {
	DB *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &repository{DB: database.DBconn}
}

func (r *repository) GetAllCategories() ([]Categoria, error) {
	var categories []Categoria
	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) GetCategory(id int) (Categoria, error) {
	var category Categoria
	if err := r.DB.First(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (r *repository) CreateCategory(category *Categoria) error {
	return r.DB.Create(category).Error
}

func (r *repository) UpdateCategory(id int, updatedCategory *Categoria) (Categoria, error) {
	var category Categoria
	if err := r.DB.First(&category, id).Error; err != nil {
		return category, err
	}
	if err := r.DB.Model(&category).Updates(updatedCategory).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (r *repository) DeleteCategory(id int) error {
	return r.DB.Delete(&Categoria{}, id).Error
}

func (r *repository) SeedCategories(categories []Categoria) error {
	return r.DB.Create(&categories).Error
}
