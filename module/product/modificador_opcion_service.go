package product

// ModificadorOpcionService contiene la lógica de negocio para ModificadorOpcion
type ModificadorOpcionService struct {
	Repo ModificadorOpcionRepository
}

func NewModificadorOpcionService(repo ModificadorOpcionRepository) *ModificadorOpcionService {
	return &ModificadorOpcionService{Repo: repo}
}

func (s *ModificadorOpcionService) GetAllModificadorOpciones(modificadorID uint64) ([]ModificadorOpcion, error) {
	return s.Repo.GetAllModificadorOpciones(modificadorID)
}

func (s *ModificadorOpcionService) GetModificadorOpcion(id uint64) (ModificadorOpcion, error) {
	return s.Repo.GetModificadorOpcion(id)
}

func (s *ModificadorOpcionService) CreateModificadorOpcion(mo *ModificadorOpcion) error {
	return s.Repo.CreateModificadorOpcion(mo)
}

func (s *ModificadorOpcionService) UpdateModificadorOpcion(id uint64, mo *ModificadorOpcion) (ModificadorOpcion, error) {
	return s.Repo.UpdateModificadorOpcion(id, mo)
}

func (s *ModificadorOpcionService) DeleteModificadorOpcion(id uint64) error {
	return s.Repo.DeleteModificadorOpcion(id)
}

func (s *ModificadorOpcionService) SeedModificadorOpciones(mos []ModificadorOpcion) error {
	return s.Repo.SeedModificadorOpciones(mos)
}

func (s *ModificadorOpcionService) PatchModificadorOpcion(id uint64, fields map[string]interface{}) (ModificadorOpcion, error) {
	return s.Repo.PatchModificadorOpcion(id, fields)
}
