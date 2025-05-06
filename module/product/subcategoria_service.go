package product

// SubCategoriaService contiene la lógica de negocio para SubCategoria
// Ahora incluye la referencia a ProductoSearchService para sincronizar Redis tras actualizaciones

type SubCategoriaService struct {
	Repo SubCategoriaRepository
	ProductoSearch *ProductoSearchService
}

// NewSubCategoriaService crea una nueva instancia del servicio
func NewSubCategoriaService(repo SubCategoriaRepository, search *ProductoSearchService) *SubCategoriaService {
	return &SubCategoriaService{Repo: repo, ProductoSearch: search}
}

// GetAllSubCategorias obtiene todas las subcategorias de la empresa
func (s *SubCategoriaService) GetAllSubCategorias(empresaID uint64) ([]SubCategoria, error) {
	return s.Repo.GetAllSubCategorias(empresaID)
}

func (s *SubCategoriaService) GetSubCategoria(id uint64, empresaID uint64) (SubCategoria, error) {
	return s.Repo.GetSubCategoria(id, empresaID)
}

func (s *SubCategoriaService) CreateSubCategoria(sub *SubCategoria) error {
	return s.Repo.CreateSubCategoria(sub)
}

// UpdateSubCategoria actualiza una subcategoria existente y reindexa productos relacionados en Redis
func (s *SubCategoriaService) UpdateSubCategoria(id uint64, empresaID uint64, sub *SubCategoria) (SubCategoria, error) {
	updated, err := s.Repo.UpdateSubCategoria(id, empresaID, sub)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		_ = s.ProductoSearch.UpdateProductsByClassification("subcategoria", sub.Nombre)
	}
	return updated, nil
}

func (s *SubCategoriaService) DeleteSubCategoria(id uint64, empresaID uint64) error {
	return s.Repo.DeleteSubCategoria(id, empresaID)
}

func (s *SubCategoriaService) PatchSubCategoria(id uint64, fields map[string]interface{}) (SubCategoria, error) {
	updated, err := s.Repo.PatchSubCategoria(id, fields)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		if nombre, ok := fields["nombre"].(string); ok {
			_ = s.ProductoSearch.UpdateProductsByClassification("subcategoria", nombre)
		}
	}
	return updated, nil
}
