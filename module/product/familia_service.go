package product

// FamiliaService contiene la lógica de negocio para Familia
type FamiliaService struct {
	Repo FamiliaRepository
}

// NewFamiliaService crea una nueva instancia del servicio
func NewFamiliaService(repo FamiliaRepository) *FamiliaService {
	return &FamiliaService{Repo: repo}
}

// GetAllFamilias obtiene todas las familias para una empresa
func (s *FamiliaService) GetAllFamilias(empresaID uint64) ([]Familia, error) {
	return s.Repo.GetAllFamilias(empresaID)
}

// GetFamiliasByCartaID obtiene todas las familias para una carta específica
func (s *FamiliaService) GetFamiliasByCartaID(cartaID uint64, empresaID uint64) ([]Familia, error) {
	return s.Repo.GetFamiliasByCartaID(cartaID, empresaID)
}

// GetFamilia obtiene una familia específica
func (s *FamiliaService) GetFamilia(id uint64, empresaID uint64) (Familia, error) {
	return s.Repo.GetFamilia(id, empresaID)
}

// CreateFamilia crea una nueva familia
func (s *FamiliaService) CreateFamilia(familia *Familia) error {
	return s.Repo.CreateFamilia(familia)
}

// UpdateFamilia actualiza una familia existente
func (s *FamiliaService) UpdateFamilia(id uint64, empresaID uint64, familia *Familia) (Familia, error) {
	return s.Repo.UpdateFamilia(id, empresaID, familia)
}

// DeleteFamilia elimina una familia
func (s *FamiliaService) DeleteFamilia(id uint64, empresaID uint64) error {
	return s.Repo.DeleteFamilia(id, empresaID)
}

// SeedFamilias inicializa múltiples familias
func (s *FamiliaService) SeedFamilias(familias []Familia) error {
	return s.Repo.SeedFamilias(familias)
}

// PatchFamilia actualiza parcialmente una familia
func (s *FamiliaService) PatchFamilia(id uint64, empresaID uint64, fields map[string]interface{}) (Familia, error) {
	return s.Repo.PatchFamilia(id, empresaID, fields)
}
