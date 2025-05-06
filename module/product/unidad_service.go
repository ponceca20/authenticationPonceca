package product

// UnidadService contiene la lógica de negocio para Unidad
type UnidadService struct {
	Repo UnidadRepository
}

// NewUnidadService crea una nueva instancia del servicio
func NewUnidadService(repo UnidadRepository) *UnidadService {
	return &UnidadService{Repo: repo}
}

// GetAllUnidades obtiene todas las unidades para una empresa
func (s *UnidadService) GetAllUnidades(empresaID uint64) ([]Unidad, error) {
	return s.Repo.GetAllUnidades(empresaID)
}

// GetUnidad obtiene una unidad específica
func (s *UnidadService) GetUnidad(id uint64, empresaID uint64) (Unidad, error) {
	return s.Repo.GetUnidad(id, empresaID)
}

// CreateUnidad crea una nueva unidad
func (s *UnidadService) CreateUnidad(unidad *Unidad) error {
	return s.Repo.CreateUnidad(unidad)
}

// UpdateUnidad actualiza una unidad existente
func (s *UnidadService) UpdateUnidad(id uint64, empresaID uint64, unidad *Unidad) (Unidad, error) {
	return s.Repo.UpdateUnidad(id, empresaID, unidad)
}

// DeleteUnidad elimina una unidad
func (s *UnidadService) DeleteUnidad(id uint64, empresaID uint64) error {
	return s.Repo.DeleteUnidad(id, empresaID)
}

// SeedUnidades inicializa múltiples unidades
func (s *UnidadService) SeedUnidades(unidades []Unidad) error {
	return s.Repo.SeedUnidades(unidades)
}

// PatchUnidad actualiza parcialmente una unidad
func (s *UnidadService) PatchUnidad(id uint64, empresaID uint64, fields map[string]interface{}) (Unidad, error) {
	return s.Repo.PatchUnidad(id, empresaID, fields)
}
