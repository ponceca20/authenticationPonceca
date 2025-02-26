package users

// RolModuloService contiene la lógica de negocio para RolModulo
type RolModuloService struct {
	Repo RolModuloRepository
}

// NewRolModuloService crea una nueva instancia del servicio
func NewRolModuloService(repo RolModuloRepository) *RolModuloService {
	return &RolModuloService{Repo: repo}
}

// GetAllRolModulos obtiene todos los rolmodulos
func (s *RolModuloService) GetAllRolModulos() ([]RolModulo, error) {
	return s.Repo.GetAllRolModulos()
}

// GetRolModulo obtiene un rolmodulo por ID
func (s *RolModuloService) GetRolModulo(id uint64) (RolModulo, error) {
	return s.Repo.GetRolModulo(id)
}

// CreateRolModulo crea un nuevo rolmodulo
func (s *RolModuloService) CreateRolModulo(rolmodulo *RolModulo) error {
	return s.Repo.CreateRolModulo(rolmodulo)
}

// UpdateRolModulo actualiza un rolmodulo existente
func (s *RolModuloService) UpdateRolModulo(id uint64, rolmodulo *RolModulo) (RolModulo, error) {
	return s.Repo.UpdateRolModulo(id, rolmodulo)
}

// DeleteRolModulo elimina un rolmodulo
func (s *RolModuloService) DeleteRolModulo(id uint64) error {
	return s.Repo.DeleteRolModulo(id)
}

// SeedRolModulos inicializa múltiples rolmodulos
func (s *RolModuloService) SeedRolModulos(rolmodulos []RolModulo) error {
	return s.Repo.SeedRolModulos(rolmodulos)
}

// PatchRolModulo actualiza parcialmente un rolmodulo
func (s *RolModuloService) PatchRolModulo(id uint64, fields map[string]interface{}) (RolModulo, error) {
	return s.Repo.PatchRolModulo(id, fields)
}
