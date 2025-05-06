package product

// CartaService contiene la lógica de negocio para Carta
type CartaService struct {
	Repo CartaRepository
}

// NewCartaService crea una nueva instancia del servicio
func NewCartaService(repo CartaRepository) *CartaService {
	return &CartaService{Repo: repo}
}

// GetAllCartas obtiene todas las cartas para una empresa
func (s *CartaService) GetAllCartas(empresaID uint64) ([]Carta, error) {
	return s.Repo.GetAllCartas(empresaID)
}

// GetCarta obtiene una carta específica
func (s *CartaService) GetCarta(id uint64, empresaID uint64) (Carta, error) {
	return s.Repo.GetCarta(id, empresaID)
}

// CreateCarta crea una nueva carta
func (s *CartaService) CreateCarta(carta *Carta) error {
	return s.Repo.CreateCarta(carta)
}

// UpdateCarta actualiza una carta existente
func (s *CartaService) UpdateCarta(id uint64, empresaID uint64, carta *Carta) (Carta, error) {
	return s.Repo.UpdateCarta(id, empresaID, carta)
}

// DeleteCarta elimina una carta
func (s *CartaService) DeleteCarta(id uint64, empresaID uint64) error {
	return s.Repo.DeleteCarta(id, empresaID)
}

// SeedCartas inicializa múltiples cartas
func (s *CartaService) SeedCartas(cartas []Carta) error {
	return s.Repo.SeedCartas(cartas)
}

// PatchCarta actualiza parcialmente una carta
func (s *CartaService) PatchCarta(id uint64, empresaID uint64, fields map[string]interface{}) (Carta, error) {
	return s.Repo.PatchCarta(id, empresaID, fields)
}
