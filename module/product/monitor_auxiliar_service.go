package product

// MonitorAuxiliarService contiene la lógica de negocio para MonitorAuxiliar
type MonitorAuxiliarService struct {
	Repo MonitorAuxiliarRepository
}

// NewMonitorAuxiliarService crea una nueva instancia del servicio
func NewMonitorAuxiliarService(repo MonitorAuxiliarRepository) *MonitorAuxiliarService {
	return &MonitorAuxiliarService{Repo: repo}
}

// GetAllMonitoresAuxiliar obtiene todos los monitores auxiliares para una empresa
func (s *MonitorAuxiliarService) GetAllMonitoresAuxiliar(empresaID uint64) ([]MonitorAuxiliar, error) {
	return s.Repo.GetAllMonitoresAuxiliar(empresaID)
}

// GetMonitorAuxiliar obtiene un monitor auxiliar específico
func (s *MonitorAuxiliarService) GetMonitorAuxiliar(id uint64, empresaID uint64) (MonitorAuxiliar, error) {
	return s.Repo.GetMonitorAuxiliar(id, empresaID)
}

// CreateMonitorAuxiliar crea un nuevo monitor auxiliar
func (s *MonitorAuxiliarService) CreateMonitorAuxiliar(monitor *MonitorAuxiliar) error {
	return s.Repo.CreateMonitorAuxiliar(monitor)
}

// UpdateMonitorAuxiliar actualiza un monitor auxiliar existente
func (s *MonitorAuxiliarService) UpdateMonitorAuxiliar(id uint64, empresaID uint64, monitor *MonitorAuxiliar) (MonitorAuxiliar, error) {
	return s.Repo.UpdateMonitorAuxiliar(id, empresaID, monitor)
}

// DeleteMonitorAuxiliar elimina un monitor auxiliar
func (s *MonitorAuxiliarService) DeleteMonitorAuxiliar(id uint64, empresaID uint64) error {
	return s.Repo.DeleteMonitorAuxiliar(id, empresaID)
}

// SeedMonitoresAuxiliar inicializa múltiples monitores auxiliares
func (s *MonitorAuxiliarService) SeedMonitoresAuxiliar(monitores []MonitorAuxiliar) error {
	return s.Repo.SeedMonitoresAuxiliar(monitores)
}

// PatchMonitorAuxiliar actualiza parcialmente un monitor auxiliar
func (s *MonitorAuxiliarService) PatchMonitorAuxiliar(id uint64, empresaID uint64, fields map[string]interface{}) (MonitorAuxiliar, error) {
	return s.Repo.PatchMonitorAuxiliar(id, empresaID, fields)
}
