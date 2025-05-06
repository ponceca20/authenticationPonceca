package product

// ProductoMediaService contiene la lógica de negocio para ProductoMedia
type ProductoMediaService struct {
	Repo ProductoMediaRepository
}

// NewProductoMediaService crea una nueva instancia del servicio
func NewProductoMediaService(repo ProductoMediaRepository) *ProductoMediaService {
	return &ProductoMediaService{Repo: repo}
}

// GetAllProductoMedias obtiene todos los medios asociados a productos para una empresa
func (s *ProductoMediaService) GetAllProductoMedias(empresaID uint64) ([]ProductoMedia, error) {
	return s.Repo.GetAllProductoMedias(empresaID)
}

// GetProductoMediasByProductoID obtiene todos los medios para un producto específico
func (s *ProductoMediaService) GetProductoMediasByProductoID(productoID uint64, empresaID uint64) ([]ProductoMedia, error) {
	return s.Repo.GetProductoMediasByProductoID(productoID, empresaID)
}

// GetProductoMedia obtiene un medio específico
func (s *ProductoMediaService) GetProductoMedia(id uint64, empresaID uint64) (ProductoMedia, error) {
	return s.Repo.GetProductoMedia(id, empresaID)
}

// CreateProductoMedia crea un nuevo medio para un producto
func (s *ProductoMediaService) CreateProductoMedia(productoMedia *ProductoMedia) error {
	return s.Repo.CreateProductoMedia(productoMedia)
}

// UpdateProductoMedia actualiza un medio existente
func (s *ProductoMediaService) UpdateProductoMedia(id uint64, empresaID uint64, productoMedia *ProductoMedia) (ProductoMedia, error) {
	return s.Repo.UpdateProductoMedia(id, empresaID, productoMedia)
}

// DeleteProductoMedia elimina un medio
func (s *ProductoMediaService) DeleteProductoMedia(id uint64, empresaID uint64) error {
	return s.Repo.DeleteProductoMedia(id, empresaID)
}

// SeedProductoMedias inicializa múltiples medios para productos
func (s *ProductoMediaService) SeedProductoMedias(productoMedias []ProductoMedia) error {
	return s.Repo.SeedProductoMedias(productoMedias)
}

// PatchProductoMedia actualiza parcialmente un medio
func (s *ProductoMediaService) PatchProductoMedia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoMedia, error) {
	return s.Repo.PatchProductoMedia(id, empresaID, fields)
}

// SetMainImage establece una imagen como principal para un producto
func (s *ProductoMediaService) SetMainImage(id uint64, productoID uint64, empresaID uint64) error {
	return s.Repo.SetMainImage(id, productoID, empresaID)
}

// ReorderMedia reordena los medios de un producto
func (s *ProductoMediaService) ReorderMedia(productoID uint64, empresaID uint64, mediaIDs []uint64) error {
	return s.Repo.ReorderMedia(productoID, empresaID, mediaIDs)
}
