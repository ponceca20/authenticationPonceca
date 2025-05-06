package product

// AreaDespachoService contiene la lógica de negocio para AreaDespacho
type AreaDespachoService struct {
	Repo AreaDespachoRepository
}

// NewAreaDespachoService crea una nueva instancia del servicio
func NewAreaDespachoService(repo AreaDespachoRepository) *AreaDespachoService {
	return &AreaDespachoService{Repo: repo}
}

// GetAllAreasDespacho obtiene todas las áreas de despacho para una empresa
func (s *AreaDespachoService) GetAllAreasDespacho(empresaID uint64) ([]AreaDespacho, error) {
	return s.Repo.GetAllAreasDespacho(empresaID)
}

// GetAreaDespacho obtiene un área de despacho específica
func (s *AreaDespachoService) GetAreaDespacho(id uint64, empresaID uint64) (AreaDespacho, error) {
	return s.Repo.GetAreaDespacho(id, empresaID)
}

// CreateAreaDespacho crea una nueva área de despacho
func (s *AreaDespachoService) CreateAreaDespacho(areaDespacho *AreaDespacho) error {
	return s.Repo.CreateAreaDespacho(areaDespacho)
}

// UpdateAreaDespacho actualiza un área de despacho existente
func (s *AreaDespachoService) UpdateAreaDespacho(id uint64, empresaID uint64, areaDespacho *AreaDespacho) (AreaDespacho, error) {
	return s.Repo.UpdateAreaDespacho(id, empresaID, areaDespacho)
}

// DeleteAreaDespacho elimina un área de despacho
func (s *AreaDespachoService) DeleteAreaDespacho(id uint64, empresaID uint64) error {
	return s.Repo.DeleteAreaDespacho(id, empresaID)
}

// SeedAreasDespacho inicializa múltiples áreas de despacho
func (s *AreaDespachoService) SeedAreasDespacho(areasDespacho []AreaDespacho) error {
	return s.Repo.SeedAreasDespacho(areasDespacho)
}

// PatchAreaDespacho actualiza parcialmente un área de despacho
func (s *AreaDespachoService) PatchAreaDespacho(id uint64, empresaID uint64, fields map[string]interface{}) (AreaDespacho, error) {
	return s.Repo.PatchAreaDespacho(id, empresaID, fields)
}
