package product

// PrecioDescuentoService contiene la lógica de negocio para PrecioDescuentoPorCantidad.
type PrecioDescuentoService struct {
	Repo PrecioDescuentoRepository
}

func NewPrecioDescuentoService(repo PrecioDescuentoRepository) *PrecioDescuentoService {
	return &PrecioDescuentoService{Repo: repo}
}

func (s *PrecioDescuentoService) GetAllPrecioDescuento(empresaID uint64) ([]PrecioDescuentoPorCantidad, error) {
	return s.Repo.GetAllPrecioDescuento(empresaID)
}

func (s *PrecioDescuentoService) GetPrecioDescuento(id uint64, empresaID uint64) (PrecioDescuentoPorCantidad, error) {
	return s.Repo.GetPrecioDescuento(id, empresaID)
}

func (s *PrecioDescuentoService) CreatePrecioDescuento(pd *PrecioDescuentoPorCantidad) error {
	return s.Repo.CreatePrecioDescuento(pd)
}

func (s *PrecioDescuentoService) UpdatePrecioDescuento(id uint64, empresaID uint64, pd *PrecioDescuentoPorCantidad) (PrecioDescuentoPorCantidad, error) {
	return s.Repo.UpdatePrecioDescuento(id, empresaID, pd)
}

func (s *PrecioDescuentoService) DeletePrecioDescuento(id uint64, empresaID uint64) error {
	return s.Repo.DeletePrecioDescuento(id, empresaID)
}

func (s *PrecioDescuentoService) PatchPrecioDescuento(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioDescuentoPorCantidad, error) {
	return s.Repo.PatchPrecioDescuento(id, empresaID, fields)
}
