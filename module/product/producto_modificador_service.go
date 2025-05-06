package product

type ProductoModificadorService struct {
	Repo ProductoModificadorRepository
}

func NewProductoModificadorService(repo ProductoModificadorRepository) *ProductoModificadorService {
	return &ProductoModificadorService{Repo: repo}
}

func (s *ProductoModificadorService) GetAllProductoModificadores(productoID uint64) ([]ProductoModificador, error) {
	if productoID == 0 {
		return s.Repo.GetAllProductoModificadores(0)
	}
	return s.Repo.GetAllProductoModificadores(productoID)
}

func (s *ProductoModificadorService) GetProductoModificador(id uint64) (ProductoModificador, error) {
	return s.Repo.GetProductoModificador(id)
}

func (s *ProductoModificadorService) CreateProductoModificador(pm *ProductoModificador) error {
	return s.Repo.CreateProductoModificador(pm)
}

func (s *ProductoModificadorService) UpdateProductoModificador(id uint64, pm *ProductoModificador) (ProductoModificador, error) {
	return s.Repo.UpdateProductoModificador(id, pm)
}

func (s *ProductoModificadorService) DeleteProductoModificador(id uint64) error {
	return s.Repo.DeleteProductoModificador(id)
}

func (s *ProductoModificadorService) SeedProductoModificadores(pms []ProductoModificador) error {
	return s.Repo.SeedProductoModificadores(pms)
}

func (s *ProductoModificadorService) PatchProductoModificador(id uint64, fields map[string]interface{}) (ProductoModificador, error) {
	return s.Repo.PatchProductoModificador(id, fields)
}
