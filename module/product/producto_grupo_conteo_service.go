package product

// ProductoGrupoConteoService contiene la lógica de negocio para ProductoGrupoConteo.
type ProductoGrupoConteoService struct {
	Repo ProductoGrupoConteoRepository
}

func NewProductoGrupoConteoService(repo ProductoGrupoConteoRepository) *ProductoGrupoConteoService {
	return &ProductoGrupoConteoService{Repo: repo}
}

func (s *ProductoGrupoConteoService) GetAllProductoGrupoConteo(empresaID uint64) ([]ProductoGrupoConteo, error) {
	return s.Repo.GetAllProductoGrupoConteo(empresaID)
}

func (s *ProductoGrupoConteoService) GetProductoGrupoConteo(id uint64, empresaID uint64) (ProductoGrupoConteo, error) {
	return s.Repo.GetProductoGrupoConteo(id, empresaID)
}

func (s *ProductoGrupoConteoService) CreateProductoGrupoConteo(pg *ProductoGrupoConteo) error {
	return s.Repo.CreateProductoGrupoConteo(pg)
}

func (s *ProductoGrupoConteoService) UpdateProductoGrupoConteo(id uint64, empresaID uint64, pg *ProductoGrupoConteo) (ProductoGrupoConteo, error) {
	return s.Repo.UpdateProductoGrupoConteo(id, empresaID, pg)
}

func (s *ProductoGrupoConteoService) DeleteProductoGrupoConteo(id uint64, empresaID uint64) error {
	return s.Repo.DeleteProductoGrupoConteo(id, empresaID)
}

func (s *ProductoGrupoConteoService) PatchProductoGrupoConteo(id uint64, empresaID uint64, fields map[string]interface{}) (ProductoGrupoConteo, error) {
	return s.Repo.PatchProductoGrupoConteo(id, empresaID, fields)
}
