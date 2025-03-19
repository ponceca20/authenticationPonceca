package imagenes

import (
	"errors"
	"practicev2/database"

	"gorm.io/gorm"
)

// ImageRepository define las operaciones de persistencia para Image
type ImageRepository interface {
	GetImageByID(id string) (Image, error)
	GetImagesByUserID(userID string) ([]Image, error)
	CreateImage(image *Image) error
	UpdateImage(image *Image) error
	DeleteImage(id string) error
}

type imageRepo struct {
	DB *gorm.DB
}

// NewImageRepository crea una nueva instancia del repositorio
func NewImageRepository() ImageRepository {
	return &imageRepo{DB: database.DBconn}
}

// GetImageByID recupera una imagen usando la clave primaria "id"
func (r *imageRepo) GetImageByID(id string) (Image, error) {
	var image Image
	err := r.DB.First(&image, id).Error // GORM usa la clave primaria por defecto
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return image, ErrImageNotFound
	}
	return image, err
}

// GetImagesByUserID recupera todas las imágenes de un usuario específico
func (r *imageRepo) GetImagesByUserID(userID string) ([]Image, error) {
	var images []Image
	err := r.DB.Where("user_id = ?", userID).Find(&images).Error
	return images, err
}

// CreateImage agrega un nuevo registro de imagen en la base de datos
func (r *imageRepo) CreateImage(image *Image) error {
	return r.DB.Create(image).Error
}

// UpdateImage actualiza el registro de una imagen (por ejemplo, para asignar rutas)
func (r *imageRepo) UpdateImage(image *Image) error {
	return r.DB.Save(image).Error
}

// DeleteImage elimina físicamente una imagen de la base de datos usando la clave primaria "id"
func (r *imageRepo) DeleteImage(id string) error {
	result := r.DB.Delete(&Image{}, id)
	if result.RowsAffected == 0 {
		return ErrImageNotFound
	}
	return result.Error
}
