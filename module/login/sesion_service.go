package users

// SesionService contiene la lógica de negocio para sesiones
type SesionService struct {
	Repo SesionRepository
}

// NewSesionService crea una nueva instancia del servicio
func NewSesionService(repo SesionRepository) *SesionService {
	return &SesionService{Repo: repo}
}

func (s *SesionService) GetAllSesions() ([]Sesion, error) {
	return s.Repo.GetAllSesions()
}

func (s *SesionService) GetSesion(id uint64) (Sesion, error) {
	return s.Repo.GetSesion(id)
}

func (s *SesionService) CreateSesion(sesion *Sesion) error {
	return s.Repo.CreateSesion(sesion)
}

func (s *SesionService) UpdateSesion(id uint64, sesion *Sesion) (Sesion, error) {
	return s.Repo.UpdateSesion(id, sesion)
}

func (s *SesionService) DeleteSesion(id uint64) error {
	return s.Repo.DeleteSesion(id)
}

func (s *SesionService) SeedSesions(sesiones []Sesion) error {
	return s.Repo.SeedSesions(sesiones)
}

func (s *SesionService) PatchSesion(id uint64, fields map[string]interface{}) (Sesion, error) {
	return s.Repo.PatchSesion(id, fields)
}
