package product

// GrupoOfertaService contiene la lógica de negocio para GrupoOferta
type GrupoOfertaService struct {
	Repo GrupoOfertaRepository
}

func NewGrupoOfertaService(repo GrupoOfertaRepository) *GrupoOfertaService {
	return &GrupoOfertaService{Repo: repo}
}

func (s *GrupoOfertaService) GetAllGrupoOfertas(empresaID uint64) ([]GrupoOferta, error) {
	return s.Repo.GetAllGrupoOfertas(empresaID)
}

func (s *GrupoOfertaService) GetGrupoOferta(id uint64, empresaID uint64) (GrupoOferta, error) {
	return s.Repo.GetGrupoOferta(id, empresaID)
}

func (s *GrupoOfertaService) CreateGrupoOferta(grupo *GrupoOferta) error {
	return s.Repo.CreateGrupoOferta(grupo)
}

func (s *GrupoOfertaService) UpdateGrupoOferta(id uint64, empresaID uint64, grupo *GrupoOferta) (GrupoOferta, error) {
	return s.Repo.UpdateGrupoOferta(id, empresaID, grupo)
}

func (s *GrupoOfertaService) DeleteGrupoOferta(id uint64, empresaID uint64) error {
	return s.Repo.DeleteGrupoOferta(id, empresaID)
}

func (s *GrupoOfertaService) SeedGrupoOfertas(grupos []GrupoOferta) error {
	return s.Repo.SeedGrupoOfertas(grupos)
}

func (s *GrupoOfertaService) PatchGrupoOferta(id uint64, empresaID uint64, fields map[string]interface{}) (GrupoOferta, error) {
	return s.Repo.PatchGrupoOferta(id, empresaID, fields)
}
