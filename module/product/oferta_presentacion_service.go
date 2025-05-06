package product

// OfertaPresentacionService contiene la lógica de negocio para OfertaPresentacion
type OfertaPresentacionService struct {
	Repo OfertaPresentacionRepository
}

func NewOfertaPresentacionService(repo OfertaPresentacionRepository) *OfertaPresentacionService {
	return &OfertaPresentacionService{Repo: repo}
}

func (s *OfertaPresentacionService) GetAllOfertaPresentaciones(empresaID uint64) ([]OfertaPresentacion, error) {
	return s.Repo.GetAllOfertaPresentaciones(empresaID)
}

func (s *OfertaPresentacionService) GetOfertaPresentacion(id uint64, empresaID uint64) (OfertaPresentacion, error) {
	return s.Repo.GetOfertaPresentacion(id, empresaID)
}

func (s *OfertaPresentacionService) CreateOfertaPresentacion(op *OfertaPresentacion) error {
	return s.Repo.CreateOfertaPresentacion(op)
}

func (s *OfertaPresentacionService) UpdateOfertaPresentacion(id uint64, empresaID uint64, op *OfertaPresentacion) (OfertaPresentacion, error) {
	return s.Repo.UpdateOfertaPresentacion(id, empresaID, op)
}

func (s *OfertaPresentacionService) DeleteOfertaPresentacion(id uint64, empresaID uint64) error {
	return s.Repo.DeleteOfertaPresentacion(id, empresaID)
}

func (s *OfertaPresentacionService) SeedOfertaPresentaciones(ops []OfertaPresentacion) error {
	return s.Repo.SeedOfertaPresentaciones(ops)
}

func (s *OfertaPresentacionService) PatchOfertaPresentacion(id uint64, empresaID uint64, fields map[string]interface{}) (OfertaPresentacion, error) {
	return s.Repo.PatchOfertaPresentacion(id, empresaID, fields)
}
