package auth

// UsuarioEmpresaService contiene la lógica de negocio para UsuarioEmpresa.
type UsuarioEmpresaService struct {
	Repo UsuarioEmpresaRepository
}

func NewUsuarioEmpresaService(repo UsuarioEmpresaRepository) *UsuarioEmpresaService {
	return &UsuarioEmpresaService{Repo: repo}
}

func (s *UsuarioEmpresaService) GetAllUsuarioEmpresas() ([]UsuarioEmpresa, error) {
	return s.Repo.GetAllUsuarioEmpresas()
}

func (s *UsuarioEmpresaService) GetUsuarioEmpresa(id uint64) (UsuarioEmpresa, error) {
	return s.Repo.GetUsuarioEmpresa(id)
}

func (s *UsuarioEmpresaService) CreateUsuarioEmpresa(ue *UsuarioEmpresa) error {
	return s.Repo.CreateUsuarioEmpresa(ue)
}

func (s *UsuarioEmpresaService) UpdateUsuarioEmpresa(id uint64, ue *UsuarioEmpresa) (UsuarioEmpresa, error) {
	return s.Repo.UpdateUsuarioEmpresa(id, ue)
}

func (s *UsuarioEmpresaService) DeleteUsuarioEmpresa(id uint64) error {
	return s.Repo.DeleteUsuarioEmpresa(id)
}

func (s *UsuarioEmpresaService) SeedUsuarioEmpresas(ues []UsuarioEmpresa) error {
	return s.Repo.SeedUsuarioEmpresas(ues)
}

func (s *UsuarioEmpresaService) PatchUsuarioEmpresa(id uint64, fields map[string]interface{}) (UsuarioEmpresa, error) {
	return s.Repo.PatchUsuarioEmpresa(id, fields)
}
