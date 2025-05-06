package product

// PrecioHistoricoService contiene la lógica de negocio para PrecioHistorico.
type PrecioHistoricoService struct {
	Repo PrecioHistoricoRepository
}

func NewPrecioHistoricoService(repo PrecioHistoricoRepository) *PrecioHistoricoService {
	return &PrecioHistoricoService{Repo: repo}
}

func (s *PrecioHistoricoService) GetAllPrecioHistorico(empresaID uint64) ([]PrecioHistorico, error) {
	return s.Repo.GetAllPrecioHistorico(empresaID)
}

func (s *PrecioHistoricoService) GetPrecioHistorico(id uint64, empresaID uint64) (PrecioHistorico, error) {
	return s.Repo.GetPrecioHistorico(id, empresaID)
}

func (s *PrecioHistoricoService) CreatePrecioHistorico(ph *PrecioHistorico) error {
	return s.Repo.CreatePrecioHistorico(ph)
}

func (s *PrecioHistoricoService) UpdatePrecioHistorico(id uint64, empresaID uint64, ph *PrecioHistorico) (PrecioHistorico, error) {
	return s.Repo.UpdatePrecioHistorico(id, empresaID, ph)
}

func (s *PrecioHistoricoService) DeletePrecioHistorico(id uint64, empresaID uint64) error {
	return s.Repo.DeletePrecioHistorico(id, empresaID)
}

func (s *PrecioHistoricoService) PatchPrecioHistorico(id uint64, empresaID uint64, fields map[string]interface{}) (PrecioHistorico, error) {
	return s.Repo.PatchPrecioHistorico(id, empresaID, fields)
}
