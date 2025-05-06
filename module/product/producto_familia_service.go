package product

// ProductoFamiliaService contiene la lógica de negocio para ProductoFamilia
type ProductoFamiliaService struct {
	Repo ProductoFamiliaRepository
}

// NewProductoFamiliaService crea una nueva instancia del servicio
func NewProductoFamiliaService(repo ProductoFamiliaRepository) *ProductoFamiliaService {
	return &ProductoFamiliaService{Repo: repo}
}

// GetAllProductoFamilias obtiene todas las relaciones ProductoFamilia para una empresa
func (s *ProductoFamiliaService) GetAllProductoFamilias(empresaID uint64) ([]ProductoFamilia, error) {
	return s.Repo.GetAllProductoFamilias(empresaID)
}

// GetProductoFamiliasByFamiliaID obtiene todas las relaciones ProductoFamilia para una familia específica
func (s *ProductoFamiliaService) GetProductoFamiliasByFamiliaID(familiaID uint64, empresaID uint64) ([]ProductoFamilia, error) {
	return s.Repo.GetProductoFamiliasByFamiliaID(familiaID, empresaID)
}

// GetProductoFamiliasByProductoID obtiene todas las relaciones ProductoFamilia para un producto específico
func (s *ProductoFamiliaService) GetProductoFamiliasByProductoID(productoID uint64, empresaID uint64) ([]ProductoFamilia, error) {
	return s.Repo.GetProductoFamiliasByProductoID(productoID, empresaID)
}

// GetProductoFamilia obtiene una relación ProductoFamilia específica
func (s *ProductoFamiliaService) GetProductoFamilia(id uint64, empresaID uint64) (ProductoFamilia, error) {
	return s.Repo.GetProductoFamilia(id, empresaID)
}

// CreateProductoFamilia crea una nueva relación ProductoFamilia
func (s *ProductoFamiliaService) CreateProductoFamilia(productoFamilia *ProductoFamilia) error {
	return s.Repo.CreateProductoFamilia(productoFamilia)
}

// UpdateProductoFamilia actualiza una relación ProductoFamilia existente
func (s *ProductoFamiliaService) UpdateProductoFamilia(id uint64, empresaID uint64, productoFamilia *ProductoFamilia) (ProductoFamilia, error) {
	return s.Repo.UpdateProductoFamilia(id, empresaID, productoFamilia)
}

// DeleteProductoFamilia elimina una relación ProductoFamilia
func (s *ProductoFamiliaService) DeleteProductoFamilia(id uint64, empresaID uint64) error {
	return s.Repo.DeleteProductoFamilia(id, empresaID)
}

// SeedProductoFamilias inicializa múltiples relaciones ProductoFamilia
func (s *ProductoFamiliaService) SeedProductoFamilias(productoFamilias []ProductoFamilia) error {
	return s.Repo.SeedProductoFamilias(productoFamilias)
}

// PatchProductoFamilia actualiza parcialmente una relación ProductoFamilia
func (s *ProductoFamiliaService) PatchProductoFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoFamilia, error) {
	return s.Repo.PatchProductoFamilia(id, empresaID, fields)
}
