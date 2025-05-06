package product

// PresentacionService contiene la lógica de negocio para Presentaciones
type PresentacionService struct {
	Repo PresentacionRepository
}

// NewPresentacionService crea una nueva instancia del servicio
func NewPresentacionService(repo PresentacionRepository) *PresentacionService {
	return &PresentacionService{Repo: repo}
}

// GetAllPresentaciones obtiene todas las presentaciones para un producto
func (s *PresentacionService) GetAllPresentaciones(productoID uint64) ([]Presentacion, error) {
	return s.Repo.GetAllPresentaciones(productoID)
}

// GetPresentacion obtiene una presentación específica
func (s *PresentacionService) GetPresentacion(id uint64) (Presentacion, error) {
	return s.Repo.GetPresentacion(id)
}

// CreatePresentacion crea una nueva presentación
func (s *PresentacionService) CreatePresentacion(presentacion *Presentacion) error {
	return s.Repo.CreatePresentacion(presentacion)
}

// UpdatePresentacion actualiza una presentación existente
func (s *PresentacionService) UpdatePresentacion(id uint64, presentacion *Presentacion) (Presentacion, error) {
	return s.Repo.UpdatePresentacion(id, presentacion)
}

// DeletePresentacion elimina una presentación
func (s *PresentacionService) DeletePresentacion(id uint64) error {
	return s.Repo.DeletePresentacion(id)
}

// SeedPresentaciones inicializa múltiples presentaciones
func (s *PresentacionService) SeedPresentaciones(presentaciones []Presentacion) error {
	return s.Repo.SeedPresentaciones(presentaciones)
}

// PatchPresentacion actualiza parcialmente una presentación
func (s *PresentacionService) PatchPresentacion(id uint64, fields map[string]interface{}) (Presentacion, error) {
	return s.Repo.PatchPresentacion(id, fields)
}

// GetAllPresentacionesByEmpresa obtiene todas las presentaciones filtradas por empresa_id
func (s *PresentacionService) GetAllPresentacionesByEmpresa(empresaID uint64) ([]Presentacion, error) {
	return s.Repo.GetAllPresentacionesByEmpresa(empresaID)
}
