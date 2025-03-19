package imagenes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ImageHandler maneja las solicitudes HTTP relacionadas con imágenes
type ImageHandler struct {
	Service *ImageService
}

// NewImageHandler crea una nueva instancia del manejador de imágenes
func NewImageHandler(s *ImageService) *ImageHandler {
	return &ImageHandler{Service: s}
}

// Métodos auxiliares para evitar duplicación

func (h *ImageHandler) getUserID(c *fiber.Ctx) (string, error) {
	// Extraer y validar token de usuario
	claims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		log.Println("Falta o es inválido el token de usuario")
		return "", fmt.Errorf("autenticación requerida")
	}
	userIDFloat, ok := claims["usuario_id"].(float64)
	if !ok {
		log.Printf("Tipo inválido para usuario_id: %T, %#v", c.Locals("user"), c.Locals("user"))
		return "", fmt.Errorf("autenticación requerida")
	}
	return fmt.Sprintf("%.0f", userIDFloat), nil
}

func (h *ImageHandler) imageMap(image *Image) fiber.Map {
	// Construir el mapa de respuesta para la imagen
	return fiber.Map{
		"id":                image.ID,
		"original_filename": image.OriginalFilename,
		"mime_type":         image.MimeType,
		"file_size":         image.FileSize,
		"image_type":        image.ImageType,
		"empresa_id":        image.EmpresaID,
		"family":            image.Family,
		// Se desreferencia image para cumplir con el tipo requerido
		"url":           h.Service.GetImageURL(*image),
		"thumbnail_url": h.Service.BaseURL + "/" + filepath.Base(image.ThumbnailPath),
		"medium_url":    h.Service.BaseURL + "/" + filepath.Base(image.MediumPath),
	}
}

// UploadImage maneja la carga de imágenes.
// Se eliminó el log de depuración y se agregan comentarios en español.
func (h *ImageHandler) UploadImage(c *fiber.Ctx) error {
	// Obtener userID mediante el método auxiliar
	userID, err := h.getUserID(c)
	if err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, "Autenticación requerida")
	}

	// Obtener y validar parámetros: image type, empresa_id y family
	imageTypeStr := c.Query("type", "user")
	imageType := ImageType(imageTypeStr)
	if !IsValidImageType(imageType) {
		return RespondWithError(c, fiber.StatusBadRequest, "Tipo de imagen inválido. Debe ser 'user' o 'product'")
	}

	empresaID := c.FormValue("empresa_id")
	if empresaID == "" {
		return RespondWithError(c, fiber.StatusBadRequest, "Falta el campo empresa_id")
	}
	family := c.FormValue("family")
	if family == "" {
		return RespondWithError(c, fiber.StatusBadRequest, "Falta el campo family")
	}

	// Obtener el archivo de imagen subido
	file, err := c.FormFile("image")
	if err != nil {
		return RespondWithError(c, fiber.StatusBadRequest, "Error al obtener el archivo subido")
	}

	// Procesar la carga de la imagen pasando los nuevos parámetros
	image, err := h.Service.UploadImage(file, userID, imageType, empresaID, family)
	if err != nil {
		status := fiber.StatusInternalServerError
		switch err {
		case ErrFileTooLarge:
			status = fiber.StatusBadRequest
		case ErrUnsupportedType:
			status = fiber.StatusUnsupportedMediaType
		}
		return RespondWithError(c, status, err.Error())
	}

	// Respuesta actualizada conforme al nuevo model
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   h.imageMap(image),
	})
}

// GetImage obtiene los detalles de una imagen específica
func (h *ImageHandler) GetImage(c *fiber.Ctx) error {
	id := c.Params("id")

	image, err := h.Service.GetImage(id)
	if err != nil {
		return RespondWithError(c, fiber.StatusNotFound, "Imagen no encontrada")
	}

	// Nota: Cada capa (handler, service, repository) tiene su función definida por separación de responsabilidades.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		// Actualización: se pasa la dirección de image para cumplir con el tipo *Image esperado
		"data": h.imageMap(&image),
	})
}

// ServeImage sirve el archivo de imagen real
func (h *ImageHandler) ServeImage(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := filepath.Join(h.Service.BasePath, filename)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return RespondWithError(c, fiber.StatusNotFound, "Archivo de imagen no encontrado")
	}

	// Configurar cabeceras de cache (1 día)
	c.Set("Cache-Control", "public, max-age=86400")

	return c.SendFile(filePath, false)
}

// GetUserImages obtiene todas las imágenes asociadas a un usuario
func (h *ImageHandler) GetUserImages(c *fiber.Ctx) error {
	userID := c.Params("userId")

	images, err := h.Service.GetImagesByUser(userID)
	if err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "Error al obtener imágenes del usuario")
	}

	// Formatear respuesta usando claves actualizadas conforme al nuevo model
	formattedImages := make([]fiber.Map, len(images))
	for i := range images {
		formattedImages[i] = h.imageMap(&images[i])
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   formattedImages,
		"count":  len(images),
	})
}

// DeleteImage elimina una imagen.
// Se eliminan los logs de depuración y se agregan comentarios en español para explicar cada paso.
func (h *ImageHandler) DeleteImage(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validar token mediante el método auxiliar
	if _, err := h.getUserID(c); err != nil {
		return RespondWithError(c, fiber.StatusUnauthorized, "Autenticación requerida")
	}

	if err := h.Service.DeleteImage(id); err != nil {
		return RespondWithError(c, fiber.StatusInternalServerError, "Error al eliminar la imagen")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Imagen eliminada exitosamente",
	})
}
