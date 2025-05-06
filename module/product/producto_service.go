package product

// ProductoService contiene la lógica de negocio para Producto

type ProductoService struct {
	Repo          ProductoRepository
	SearchService *ProductoSearchService
}

// NewProductoService crea una nueva instancia del servicio
func NewProductoService(repo ProductoRepository, search *ProductoSearchService) *ProductoService {
	return &ProductoService{Repo: repo, SearchService: search}
}

// GetAllProductos obtiene todos los productos de una empresa
func (s *ProductoService) GetAllProductos(empresaID uint64) ([]Producto, error) {
	return s.Repo.GetAllProductos(empresaID)
}

// GetProducto obtiene un producto específico
func (s *ProductoService) GetProducto(id uint64, empresaID uint64) (Producto, error) {
	return s.Repo.GetProducto(id, empresaID)
}

// CreateProducto crea un nuevo producto y lo indexa en Redis
func (s *ProductoService) CreateProducto(p *Producto) error {
	err := s.Repo.CreateProducto(p)
	if err != nil {
		return err
	}
	if s.SearchService != nil {
		_ = s.SearchService.IndexProductByID(uint64(p.ID))
	}
	return nil
}

// UpdateProducto actualiza un producto existente y lo reindexa en Redis
func (s *ProductoService) UpdateProducto(id uint64, empresaID uint64, p *Producto) (Producto, error) {
	prod, err := s.Repo.UpdateProducto(id, empresaID, p)
	if err != nil {
		return prod, err
	}
	if s.SearchService != nil {
		_ = s.SearchService.IndexProductByID(id)
	}
	return prod, nil
}

// DeleteProducto elimina un producto y lo elimina de Redis
func (s *ProductoService) DeleteProducto(id uint64, empresaID uint64) error {
	err := s.Repo.DeleteProducto(id, empresaID)
	if err != nil {
		return err
	}
	if s.SearchService != nil {
		_ = s.SearchService.DeleteProductByID(id)
	}
	return nil
}

// PatchProducto actualiza parcialmente un producto y lo reindexa en Redis
func (s *ProductoService) PatchProducto(id uint64, empresaID uint64, fields map[string]interface{}) (Producto, error) {
	prod, err := s.Repo.PatchProducto(id, empresaID, fields)
	if err != nil {
		return prod, err
	}
	if s.SearchService != nil {
		_ = s.SearchService.IndexProductByID(id)
	}
	return prod, nil
}
