package product

// PresentacionCompuestaService contiene la lógica de negocio para PresentacionCompuesta
type PresentacionCompuestaService struct {
	Repo PresentacionCompuestaRepository
}

func NewPresentacionCompuestaService(repo PresentacionCompuestaRepository) *PresentacionCompuestaService {
	return &PresentacionCompuestaService{Repo: repo}
}

func (s *PresentacionCompuestaService) GetAllPresentacionCompuestas(empresaID uint64) ([]PresentacionCompuesta, error) {
	return s.Repo.GetAllPresentacionCompuestas(empresaID)
}

func (s *PresentacionCompuestaService) GetPresentacionCompuesta(id uint64, empresaID uint64) (PresentacionCompuesta, error) {
	return s.Repo.GetPresentacionCompuesta(id, empresaID)
}

func (s *PresentacionCompuestaService) CreatePresentacionCompuesta(pc *PresentacionCompuesta) error {
	return s.Repo.CreatePresentacionCompuesta(pc)
}

func (s *PresentacionCompuestaService) UpdatePresentacionCompuesta(id uint64, empresaID uint64, pc *PresentacionCompuesta) (PresentacionCompuesta, error) {
	return s.Repo.UpdatePresentacionCompuesta(id, empresaID, pc)
}

func (s *PresentacionCompuestaService) DeletePresentacionCompuesta(id uint64, empresaID uint64) error {
	return s.Repo.DeletePresentacionCompuesta(id, empresaID)
}

func (s *PresentacionCompuestaService) SeedPresentacionCompuestas(pcs []PresentacionCompuesta) error {
	return s.Repo.SeedPresentacionCompuestas(pcs)
}

func (s *PresentacionCompuestaService) PatchPresentacionCompuesta(id uint64, empresaID uint64, fields map[string]interface{}) (PresentacionCompuesta, error) {
	return s.Repo.PatchPresentacionCompuesta(id, empresaID, fields)
}
