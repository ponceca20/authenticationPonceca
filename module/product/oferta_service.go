package product

// OfertaService contiene la lógica de negocio para Oferta
type OfertaService struct {
	Repo OfertaRepository
}

// NewOfertaService crea una nueva instancia del servicio de Oferta
func NewOfertaService(repo OfertaRepository) *OfertaService {
	return &OfertaService{Repo: repo}
}

func (s *OfertaService) GetAllOfertas(empresaID uint64) ([]Oferta, error) {
	return s.Repo.GetAllOfertas(empresaID)
}

func (s *OfertaService) GetOferta(id uint64, empresaID uint64) (Oferta, error) {
	return s.Repo.GetOferta(id, empresaID)
}

func (s *OfertaService) CreateOferta(oferta *Oferta) error {
	return s.Repo.CreateOferta(oferta)
}

func (s *OfertaService) UpdateOferta(id uint64, empresaID uint64, oferta *Oferta) (Oferta, error) {
	return s.Repo.UpdateOferta(id, empresaID, oferta)
}

func (s *OfertaService) DeleteOferta(id uint64, empresaID uint64) error {
	return s.Repo.DeleteOferta(id, empresaID)
}

func (s *OfertaService) SeedOfertas(ofertas []Oferta) error {
	return s.Repo.SeedOfertas(ofertas)
}

func (s *OfertaService) PatchOferta(id uint64, empresaID uint64, fields map[string]interface{}) (Oferta, error) {
	return s.Repo.PatchOferta(id, empresaID, fields)
}
