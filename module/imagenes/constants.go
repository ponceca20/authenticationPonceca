package imagenes

// Constants for image handling
const (
	MaxFileSize = 30 * 1024 * 1024 // 10MB
)

// Allowed MIME types
var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	// Nuevos formatos añadidos:
	"image/gif": true,
	"image/bmp": true,
}
