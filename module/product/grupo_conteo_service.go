package product

// GrupoConteoService contiene la lógica de negocio para GrupoConteo.
type GrupoConteoService struct {
	Repo GrupoConteoRepository
}

// NewGrupoConteoService crea una nueva instancia del servicio.
func NewGrupoConteoService(repo GrupoConteoRepository) *GrupoConteoService {
	return &GrupoConteoService{Repo: repo}
}

func (s *GrupoConteoService) GetAll(empresaID uint64) ([]GrupoConteo, error) {
	return s.Repo.GetAll(empresaID)
}

func (s *GrupoConteoService) GetByID(id uint64, empresaID uint64) (GrupoConteo, error) {
	return s.Repo.GetByID(id, empresaID)
}

func (s *GrupoConteoService) Create(item *GrupoConteo) error {
	return s.Repo.Create(item)
}

func (s *GrupoConteoService) Update(id uint64, empresaID uint64, item *GrupoConteo) (GrupoConteo, error) {
	return s.Repo.Update(id, empresaID, item)
}

func (s *GrupoConteoService) Delete(id uint64, empresaID uint64) error {
	return s.Repo.Delete(id, empresaID)
}

func (s *GrupoConteoService) Patch(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoConteo, error) {
	return s.Repo.Patch(id, empresaID, fields)
}
