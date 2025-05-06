package product

// PresentacionPesoService contiene la lógica de negocio para PresentacionPeso
type PresentacionPesoService struct {
	Repo PresentacionPesoRepository
}

// NewPresentacionPesoService crea una nueva instancia del servicio
func NewPresentacionPesoService(repo PresentacionPesoRepository) *PresentacionPesoService {
	return &PresentacionPesoService{Repo: repo}
}

// GetAllPresentacionesPeso obtiene todas las configuraciones de peso para una presentación
func (s *PresentacionPesoService) GetAllPresentacionesPeso(presentacionID uint64) ([]PresentacionPeso, error) {
	return s.Repo.GetAllPresentacionesPeso(presentacionID)
}

// GetPresentacionPeso obtiene una configuración de peso específica
func (s *PresentacionPesoService) GetPresentacionPeso(id uint64) (PresentacionPeso, error) {
	return s.Repo.GetPresentacionPeso(id)
}

// CreatePresentacionPeso crea una nueva configuración de peso
func (s *PresentacionPesoService) CreatePresentacionPeso(peso *PresentacionPeso) error {
	return s.Repo.CreatePresentacionPeso(peso)
}

// UpdatePresentacionPeso actualiza una configuración de peso existente
func (s *PresentacionPesoService) UpdatePresentacionPeso(id uint64, peso *PresentacionPeso) (PresentacionPeso, error) {
	return s.Repo.UpdatePresentacionPeso(id, peso)
}

// DeletePresentacionPeso elimina una configuración de peso
func (s *PresentacionPesoService) DeletePresentacionPeso(id uint64) error {
	return s.Repo.DeletePresentacionPeso(id)
}

// SeedPresentacionesPeso inicializa múltiples configuraciones de peso
func (s *PresentacionPesoService) SeedPresentacionesPeso(pesos []PresentacionPeso) error {
	return s.Repo.SeedPresentacionesPeso(pesos)
}

// PatchPresentacionPeso actualiza parcialmente una configuración de peso
func (s *PresentacionPesoService) PatchPresentacionPeso(id uint64, fields map[string]interface{}) (PresentacionPeso, error) {
	return s.Repo.PatchPresentacionPeso(id, fields)
}
