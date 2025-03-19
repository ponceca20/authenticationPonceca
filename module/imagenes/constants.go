package imagenes

// Constants for image handling
const (
	MaxFileSize = 60 * 1024 * 1024 // 5MB
)

// Allowed MIME types
var allowedMimeTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/webp":    true,
	"image/svg+xml": true,
}
