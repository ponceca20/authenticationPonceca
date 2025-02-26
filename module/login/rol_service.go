package users

// RolService contiene la lógica de negocio para Rol.
type RolService struct {
	Repo RolRepository
}

// NewRolService crea una nueva instancia del servicio.
func NewRolService(repo RolRepository) *RolService {
	return &RolService{Repo: repo}
}

func (s *RolService) GetAllRoles() ([]Rol, error) {
	return s.Repo.GetAllRoles()
}

func (s *RolService) GetRol(id uint64) (Rol, error) {
	return s.Repo.GetRol(id)
}

func (s *RolService) CreateRol(rol *Rol) error {
	return s.Repo.CreateRol(rol)
}

func (s *RolService) UpdateRol(id uint64, rol *Rol) (Rol, error) {
	return s.Repo.UpdateRol(id, rol)
}

func (s *RolService) DeleteRol(id uint64) error {
	return s.Repo.DeleteRol(id)
}

func (s *RolService) SeedRoles(roles []Rol) error {
	return s.Repo.SeedRoles(roles)
}

func (s *RolService) PatchRol(id uint64, fields map[string]interface{}) (Rol, error) {
	return s.Repo.PatchRol(id, fields)
}
