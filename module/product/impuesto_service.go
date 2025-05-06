package product

// ImpuestoService contiene la lógica de negocio para Impuesto
type ImpuestoService struct {
	Repo ImpuestoRepository
}

func NewImpuestoService(repo ImpuestoRepository) *ImpuestoService {
	return &ImpuestoService{Repo: repo}
}

func (s *ImpuestoService) GetAllImpuestos(empresaID uint64) ([]Impuesto, error) {
	return s.Repo.GetAllImpuestos(empresaID)
}

func (s *ImpuestoService) GetImpuesto(id uint64, empresaID uint64) (Impuesto, error) {
	return s.Repo.GetImpuesto(id, empresaID)
}

func (s *ImpuestoService) CreateImpuesto(impuesto *Impuesto) error {
	return s.Repo.CreateImpuesto(impuesto)
}

func (s *ImpuestoService) UpdateImpuesto(id uint64, empresaID uint64, impuesto *Impuesto) (Impuesto, error) {
	return s.Repo.UpdateImpuesto(id, empresaID, impuesto)
}

func (s *ImpuestoService) DeleteImpuesto(id uint64, empresaID uint64) error {
	return s.Repo.DeleteImpuesto(id, empresaID)
}

func (s *ImpuestoService) SeedImpuestos(impuestos []Impuesto) error {
	return s.Repo.SeedImpuestos(impuestos)
}

func (s *ImpuestoService) PatchImpuesto(id uint64, empresaID uint64, fields map[string]interface{}) (Impuesto, error) {
	return s.Repo.PatchImpuesto(id, empresaID, fields)
}
