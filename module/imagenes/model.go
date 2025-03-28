package imagenes

import (
	"time"

	"gorm.io/gorm"
)

// ImageType represents the context where the image is used
type ImageType string

const (
	UserImage    ImageType = "user"
	ProductImage ImageType = "product"
)

// Image represents an image record in the database
type Image struct {
	gorm.Model
	UserID           string    `json:"user_id" gorm:"type:varchar(36);index"`
	OriginalFilename string    `json:"original_filename" gorm:"type:varchar(255)"`
	MimeType         string    `json:"mime_type" gorm:"type:varchar(100)"`
	FileSize         int       `json:"file_size" gorm:"type:int"`
	ImageType        ImageType `json:"image_type" gorm:"type:varchar(50);index"`
	FilePath         string    `json:"file_path" gorm:"type:varchar(255)"`
	ThumbnailPath    string    `json:"thumbnail_path" gorm:"type:varchar(255)"`
	MediumPath       string    `json:"medium_path" gorm:"type:varchar(255)"`
	EmpresaID        uint      `json:"empresa_id" gorm:"index"`
	Family           string    `json:"family" gorm:"type:varchar(50);index"`
	StatusPermanente bool      `json:"status_permanente" gorm:"default:false"`
	ExpiresAt        time.Time `json:"expires_at" gorm:"index"`
}

// TableName specifies the table name for the Image model
func (Image) TableName() string {
	return "images"
}
