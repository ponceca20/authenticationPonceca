package product

// TipoProductoService contiene la lógica de negocio para TipoProducto
// Ahora incluye la referencia a ProductoSearchService para sincronizar Redis tras actualizaciones

type TipoProductoService struct {
	Repo           TipoProductoRepository
	ProductoSearch *ProductoSearchService
}

// NewTipoProductoService crea una nueva instancia del servicio
func NewTipoProductoService(repo TipoProductoRepository, search *ProductoSearchService) *TipoProductoService {
	return &TipoProductoService{Repo: repo, ProductoSearch: search}
}

// GetAllTipoProductos obtiene todos los tipos de producto de una empresa
func (s *TipoProductoService) GetAllTipoProductos(empresaID uint64) ([]TipoProducto, error) {
	return s.Repo.GetAllTipoProductos(empresaID)
}

// GetTipoProducto obtiene un tipo de producto específico
func (s *TipoProductoService) GetTipoProducto(id uint64, empresaID uint64) (TipoProducto, error) {
	return s.Repo.GetTipoProducto(id, empresaID)
}

// CreateTipoProducto crea un nuevo tipo de producto
func (s *TipoProductoService) CreateTipoProducto(tp *TipoProducto) error {
	return s.Repo.CreateTipoProducto(tp)
}

// UpdateTipoProducto actualiza un tipo de producto existente y reindexa productos relacionados en Redis
func (s *TipoProductoService) UpdateTipoProducto(id uint64, empresaID uint64, tp *TipoProducto) (TipoProducto, error) {
	updated, err := s.Repo.UpdateTipoProducto(id, empresaID, tp)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		_ = s.ProductoSearch.UpdateProductsByClassification("tipo_producto", tp.Nombre)
	}
	return updated, nil
}

// DeleteTipoProducto elimina un tipo de producto
func (s *TipoProductoService) DeleteTipoProducto(id uint64, empresaID uint64) error {
	return s.Repo.DeleteTipoProducto(id, empresaID)
}

// PatchTipoProducto actualiza parcialmente un tipo de producto
func (s *TipoProductoService) PatchTipoProducto(id uint64, empresaID uint64, fields map[string]interface{}) (TipoProducto, error) {
	updated, err := s.Repo.PatchTipoProducto(id, empresaID, fields)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		if nombre, ok := fields["nombre"].(string); ok {
			_ = s.ProductoSearch.UpdateProductsByClassification("tipo_producto", nombre)
		}
	}
	return updated, nil
}
