package product

// PresentacionTiempoService contiene la lógica de negocio para PresentacionTiempo
type PresentacionTiempoService struct {
	Repo PresentacionTiempoRepository
}

// NewPresentacionTiempoService crea una nueva instancia del servicio
func NewPresentacionTiempoService(repo PresentacionTiempoRepository) *PresentacionTiempoService {
	return &PresentacionTiempoService{Repo: repo}
}

// GetAllPresentacionesTiempo obtiene todas las configuraciones de tiempo para una presentación
func (s *PresentacionTiempoService) GetAllPresentacionesTiempo(presentacionID uint64) ([]PresentacionTiempo, error) {
	return s.Repo.GetAllPresentacionesTiempo(presentacionID)
}

// GetPresentacionTiempo obtiene una configuración de tiempo específica
func (s *PresentacionTiempoService) GetPresentacionTiempo(id uint64) (PresentacionTiempo, error) {
	return s.Repo.GetPresentacionTiempo(id)
}

// CreatePresentacionTiempo crea una nueva configuración de tiempo
func (s *PresentacionTiempoService) CreatePresentacionTiempo(tiempo *PresentacionTiempo) error {
	return s.Repo.CreatePresentacionTiempo(tiempo)
}

// UpdatePresentacionTiempo actualiza una configuración de tiempo existente
func (s *PresentacionTiempoService) UpdatePresentacionTiempo(id uint64, tiempo *PresentacionTiempo) (PresentacionTiempo, error) {
	return s.Repo.UpdatePresentacionTiempo(id, tiempo)
}

// DeletePresentacionTiempo elimina una configuración de tiempo
func (s *PresentacionTiempoService) DeletePresentacionTiempo(id uint64) error {
	return s.Repo.DeletePresentacionTiempo(id)
}

// SeedPresentacionesTiempo inicializa múltiples configuraciones de tiempo
func (s *PresentacionTiempoService) SeedPresentacionesTiempo(tiempos []PresentacionTiempo) error {
	return s.Repo.SeedPresentacionesTiempo(tiempos)
}

// PatchPresentacionTiempo actualiza parcialmente una configuración de tiempo
func (s *PresentacionTiempoService) PatchPresentacionTiempo(id uint64, fields map[string]interface{}) (PresentacionTiempo, error) {
	return s.Repo.PatchPresentacionTiempo(id, fields)
}
