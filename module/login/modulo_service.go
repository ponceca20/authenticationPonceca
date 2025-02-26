package users

// ModuloService contiene la lógica de negocio para módulos
type ModuloService struct {
	Repo ModuloRepository
}

// NewModuloService crea una nueva instancia del servicio
func NewModuloService(repo ModuloRepository) *ModuloService {
	return &ModuloService{Repo: repo}
}

func (s *ModuloService) GetAllModulos() ([]Modulo, error) {
	return s.Repo.GetAllModulos()
}

func (s *ModuloService) GetModulo(id uint64) (Modulo, error) {
	return s.Repo.GetModulo(id)
}

func (s *ModuloService) CreateModulo(modulo *Modulo) error {
	return s.Repo.CreateModulo(modulo)
}

func (s *ModuloService) UpdateModulo(id uint64, modulo *Modulo) (Modulo, error) {
	return s.Repo.UpdateModulo(id, modulo)
}

func (s *ModuloService) DeleteModulo(id uint64) error {
	return s.Repo.DeleteModulo(id)
}

func (s *ModuloService) SeedModulos(modulos []Modulo) error {
	return s.Repo.SeedModulos(modulos)
}

func (s *ModuloService) PatchModulo(id uint64, fields map[string]interface{}) (Modulo, error) {
	return s.Repo.PatchModulo(id, fields)
}
