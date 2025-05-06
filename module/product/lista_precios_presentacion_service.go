package product

// ListaPreciosPresentacionService contiene la lógica de negocio para ListaPreciosPresentacion.
type ListaPreciosPresentacionService struct {
	Repo ListaPreciosPresentacionRepository
}

// NewListaPreciosPresentacionService crea una nueva instancia del servicio.
func NewListaPreciosPresentacionService(repo ListaPreciosPresentacionRepository) *ListaPreciosPresentacionService {
	return &ListaPreciosPresentacionService{Repo: repo}
}

func (s *ListaPreciosPresentacionService) GetAll() ([]ListaPreciosPresentacion, error) {
	return s.Repo.GetAll()
}

func (s *ListaPreciosPresentacionService) GetByID(id uint64) (ListaPreciosPresentacion, error) {
	return s.Repo.GetByID(id)
}

func (s *ListaPreciosPresentacionService) Create(item *ListaPreciosPresentacion) error {
	return s.Repo.Create(item)
}

func (s *ListaPreciosPresentacionService) Update(id uint64, item *ListaPreciosPresentacion) (ListaPreciosPresentacion, error) {
	return s.Repo.Update(id, item)
}

func (s *ListaPreciosPresentacionService) Delete(id uint64) error {
	return s.Repo.Delete(id)
}

func (s *ListaPreciosPresentacionService) Patch(id uint64, fields map[string]interface{}) (ListaPreciosPresentacion, error) {
	return s.Repo.Patch(id, fields)
}
