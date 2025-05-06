package product

// MonitorCocinaService contiene la lógica de negocio para MonitorCocina
type MonitorCocinaService struct {
	Repo MonitorCocinaRepository
}

// NewMonitorCocinaService crea una nueva instancia del servicio
func NewMonitorCocinaService(repo MonitorCocinaRepository) *MonitorCocinaService {
	return &MonitorCocinaService{Repo: repo}
}

// GetAllMonitoresCocina obtiene todos los monitores de cocina para una empresa
func (s *MonitorCocinaService) GetAllMonitoresCocina(empresaID uint64) ([]MonitorCocina, error) {
	return s.Repo.GetAllMonitoresCocina(empresaID)
}

// GetMonitorCocina obtiene un monitor de cocina específico
func (s *MonitorCocinaService) GetMonitorCocina(id uint64, empresaID uint64) (MonitorCocina, error) {
	return s.Repo.GetMonitorCocina(id, empresaID)
}

// CreateMonitorCocina crea un nuevo monitor de cocina
func (s *MonitorCocinaService) CreateMonitorCocina(monitor *MonitorCocina) error {
	return s.Repo.CreateMonitorCocina(monitor)
}

// UpdateMonitorCocina actualiza un monitor de cocina existente
func (s *MonitorCocinaService) UpdateMonitorCocina(id uint64, empresaID uint64, monitor *MonitorCocina) (MonitorCocina, error) {
	return s.Repo.UpdateMonitorCocina(id, empresaID, monitor)
}

// DeleteMonitorCocina elimina un monitor de cocina
func (s *MonitorCocinaService) DeleteMonitorCocina(id uint64, empresaID uint64) error {
	return s.Repo.DeleteMonitorCocina(id, empresaID)
}

// SeedMonitoresCocina inicializa múltiples monitores de cocina
func (s *MonitorCocinaService) SeedMonitoresCocina(monitores []MonitorCocina) error {
	return s.Repo.SeedMonitoresCocina(monitores)
}

// PatchMonitorCocina actualiza parcialmente un monitor de cocina
func (s *MonitorCocinaService) PatchMonitorCocina(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorCocina, error) {
	return s.Repo.PatchMonitorCocina(id, empresaID, fields)
}
