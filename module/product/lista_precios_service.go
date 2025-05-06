package product

// ListaPreciosService contiene la lógica de negocio para ListaPrecios
type ListaPreciosService struct {
	Repo ListaPreciosRepository
}

// NewListaPreciosService crea una nueva instancia del servicio
func NewListaPreciosService(repo ListaPreciosRepository) *ListaPreciosService {
	return &ListaPreciosService{Repo: repo}
}

// GetAllListaPrecios obtiene todas las listas de precios para una empresa
func (s *ListaPreciosService) GetAllListaPrecios(empresaID uint64) ([]ListaPrecios, error) {
	return s.Repo.GetAllListaPrecios(empresaID)
}

// GetListaPreciosWithPresentaciones obtiene una lista de precios con sus presentaciones asociadas
func (s *ListaPreciosService) GetListaPreciosWithPresentaciones(id uint64, empresaID uint64) (ListaPrecios, error) {
	return s.Repo.GetListaPreciosWithPresentaciones(id, empresaID)
}

// GetListaPrecios obtiene una lista de precios específica
func (s *ListaPreciosService) GetListaPrecios(id uint64, empresaID uint64) (ListaPrecios, error) {
	return s.Repo.GetListaPrecios(id, empresaID)
}

// CreateListaPrecios crea una nueva lista de precios
func (s *ListaPreciosService) CreateListaPrecios(listaPrecios *ListaPrecios) error {
	return s.Repo.CreateListaPrecios(listaPrecios)
}

// UpdateListaPrecios actualiza una lista de precios existente
func (s *ListaPreciosService) UpdateListaPrecios(id uint64, empresaID uint64, listaPrecios *ListaPrecios) (ListaPrecios, error) {
	return s.Repo.UpdateListaPrecios(id, empresaID, listaPrecios)
}

// DeleteListaPrecios elimina una lista de precios
func (s *ListaPreciosService) DeleteListaPrecios(id uint64, empresaID uint64) error {
	return s.Repo.DeleteListaPrecios(id, empresaID)
}

// SeedListaPrecios inicializa múltiples listas de precios
func (s *ListaPreciosService) SeedListaPrecios(listaPrecios []ListaPrecios) error {
	return s.Repo.SeedListaPrecios(listaPrecios)
}

// PatchListaPrecios actualiza parcialmente una lista de precios
func (s *ListaPreciosService) PatchListaPrecios(id uint64, empresaID uint64, fields map[string]interface{}) (ListaPrecios, error) {
	return s.Repo.PatchListaPrecios(id, empresaID, fields)
}
