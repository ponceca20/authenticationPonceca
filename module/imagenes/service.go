package imagenes

import (
	"errors"
	"fmt"
	"image" // agregado para usar image.Image
	"log"   // agregado para logging
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv" // nueva importación

	"github.com/disintegration/imaging"
)

// ImageService contains the business logic for images
type ImageService struct {
	Repo         ImageRepository
	BasePath     string
	BaseURL      string
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

// NewImageService creates a new instance of the image service
func NewImageService(repo ImageRepository, basePath, baseURL string) *ImageService {
	// Create directory if it doesn't exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			log.Printf("Failed to create base directory '%s': %v", basePath, err)
		}
	}

	return &ImageService{
		Repo:         repo,
		BasePath:     basePath,
		BaseURL:      baseURL,
		MaxFileSize:  MaxFileSize,
		AllowedTypes: allowedMimeTypes,
	}
}

// UploadImage maneja la carga de una nueva imagen utilizando la ID del registro
func (s *ImageService) UploadImage(file *multipart.FileHeader, userID string, imageType ImageType, empresaID, family string) (*Image, error) {
	// Validaciones
	if file.Size > s.MaxFileSize {
		return nil, ErrFileTooLarge
	}
	mimeType := file.Header.Get("Content-Type")
	if !s.AllowedTypes[mimeType] {
		return nil, ErrUnsupportedType
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = GetFileExtensionFromMime(mimeType)
	}

	// Decodificar la imagen
	src, err := file.Open()
	if err != nil {
		log.Printf("Error al abrir el archivo: %v", err)
		return nil, err
	}
	defer src.Close()
	img, err := imaging.Decode(src)
	if err != nil {
		log.Printf("Error al decodificar la imagen: %v", err)
		return nil, err
	}

	// Redimensionar la imagen original si el ancho es mayor a 1080
	if img.Bounds().Dx() > 1080 {
		img = imaging.Resize(img, 1080, 0, imaging.Lanczos)
	}

	// Convertir empresaID de string a uint
	empID, err := strconv.ParseUint(empresaID, 10, 64)
	if err != nil {
		return nil, errors.New("EmpresaID debe ser un número entero sin signo válido")
	}

	// Crear registro en la base de datos sin rutas aún
	imageRecord := &Image{
		UserID:           userID,
		OriginalFilename: file.Filename,
		MimeType:         mimeType,
		FileSize:         int(file.Size),
		ImageType:        imageType,
		EmpresaID:        uint(empID),
		Family:           family,
		// FilePath, ThumbnailPath y MediumPath se asignarán luego
	}
	if err := s.Repo.CreateImage(imageRecord); err != nil {
		return nil, err
	}

	// Usar la ID del registro para construir las rutas de los archivos
	idStr := fmt.Sprintf("%d", imageRecord.ID)
	originalPath := filepath.Join(s.BasePath, idStr+ext)
	thumbPath := filepath.Join(s.BasePath, idStr+"_thumb"+ext)
	mediumPath := filepath.Join(s.BasePath, idStr+"_medium"+ext)

	// Guardar imagen original
	if err = imaging.Save(img, originalPath); err != nil {
		log.Printf("Error al guardar imagen original: %v", err)
		return nil, err
	}

	// Redimensionar según ancho, evitando upscaling
	originalWidth := img.Bounds().Dx()
	var thumbImg, mediumImg image.Image
	if originalWidth > 150 {
		thumbImg = imaging.Resize(img, 150, 0, imaging.Lanczos)
	} else {
		thumbImg = img
	}
	if originalWidth > 500 {
		mediumImg = imaging.Resize(img, 500, 0, imaging.Lanczos)
	} else {
		mediumImg = img
	}

	if err = imaging.Save(thumbImg, thumbPath); err != nil {
		log.Printf("Error al guardar la imagen miniatura: %v", err)
		os.Remove(originalPath)
		// Se podría eliminar el registro en caso de error
		return nil, err
	}
	if err = imaging.Save(mediumImg, mediumPath); err != nil {
		log.Printf("Error al guardar la imagen mediana: %v", err)
		os.Remove(originalPath)
		os.Remove(thumbPath)
		return nil, err
	}

	// Actualizar el registro con las rutas físicas de los archivos
	imageRecord.FilePath = originalPath
	imageRecord.ThumbnailPath = thumbPath
	imageRecord.MediumPath = mediumPath
	if err := s.Repo.UpdateImage(imageRecord); err != nil {
		// Limpiar archivos si la actualización falla
		os.Remove(originalPath)
		os.Remove(thumbPath)
		os.Remove(mediumPath)
		return nil, err
	}

	return imageRecord, nil
}

// GetImage retrieves an image by ID
func (s *ImageService) GetImage(id string) (Image, error) {
	return s.Repo.GetImageByID(id)
}

// GetImagesByUser retrieves all images for a specific user
func (s *ImageService) GetImagesByUser(userID string) ([]Image, error) {
	return s.Repo.GetImagesByUserID(userID)
}

// DeleteImage deletes an image and its file versions
func (s *ImageService) DeleteImage(id string) error {
	// Get image details
	image, err := s.Repo.GetImageByID(id)
	if err != nil {
		return err
	}

	// Delete original file and other versions
	paths := []string{image.FilePath, image.ThumbnailPath, image.MediumPath}
	for _, path := range paths {
		if err := os.Remove(path); err != nil {
			log.Printf("Warning: failed to delete file '%s': %v", path, err)
		}
	}

	// Delete from database
	return s.Repo.DeleteImage(id)
}

// GetImageURL genera un URL para acceder a la imagen
func (s *ImageService) GetImageURL(image Image) string {
	filename := filepath.Base(image.FilePath)
	return fmt.Sprintf("%s/%s", s.BaseURL, filename)
}
