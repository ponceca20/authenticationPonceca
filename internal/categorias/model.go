package categorias

import "gorm.io/gorm"

type Categoria struct {
	gorm.Model
	Name string `json:"name"`
}
