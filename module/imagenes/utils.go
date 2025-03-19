package imagenes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Common errors
var (
	ErrInvalidImageType = errors.New("invalid image type")
	ErrFileTooLarge     = errors.New("file size exceeds maximum limit")
	ErrUnsupportedType  = errors.New("unsupported file type")
	ErrImageNotFound    = errors.New("image not found")
)

// RespondWithError creates a consistent error response
func RespondWithError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

// GetFileExtensionFromMime returns the appropriate file extension for a MIME type
func GetFileExtensionFromMime(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	default:
		return ""
	}
}

// IsValidImageType checks if the provided type is valid
func IsValidImageType(imageType ImageType) bool {
	switch imageType {
	case UserImage, ProductImage:
		return true
	default:
		return false
	}
}
