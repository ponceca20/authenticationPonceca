package product

// ModificadorService contiene la lógica de negocio para Modificador
type ModificadorService struct {
	Repo ModificadorRepository
}

func NewModificadorService(repo ModificadorRepository) *ModificadorService {
	return &ModificadorService{Repo: repo}
}

func (s *ModificadorService) GetAllModificadores(empresaID uint64) ([]Modificador, error) {
	return s.Repo.GetAllModificadores(empresaID)
}

func (s *ModificadorService) GetModificador(id uint64, empresaID uint64) (Modificador, error) {
	return s.Repo.GetModificador(id, empresaID)
}

func (s *ModificadorService) CreateModificador(m *Modificador) error {
	return s.Repo.CreateModificador(m)
}

func (s *ModificadorService) UpdateModificador(id uint64, empresaID uint64, m *Modificador) (Modificador, error) {
	return s.Repo.UpdateModificador(id, empresaID, m)
}

func (s *ModificadorService) DeleteModificador(id uint64, empresaID uint64) error {
	return s.Repo.DeleteModificador(id, empresaID)
}

func (s *ModificadorService) SeedModificadores(mods []Modificador) error {
	return s.Repo.SeedModificadores(mods)
}

func (s *ModificadorService) PatchModificador(id uint64, empresaID uint64, fields map[string]interface{}) (Modificador, error) {
	return s.Repo.PatchModificador(id, empresaID, fields)
}
