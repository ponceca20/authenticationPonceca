package product

// ClaseService contiene la lógica de negocio para Clase
// Ahora incluye la referencia a ProductoSearchService para sincronizar Redis tras actualizaciones

type ClaseService struct {
	Repo ClaseRepository
	ProductoSearch *ProductoSearchService
}

// NewClaseService crea una nueva instancia del servicio
func NewClaseService(repo ClaseRepository, search *ProductoSearchService) *ClaseService {
	return &ClaseService{Repo: repo, ProductoSearch: search}
}

// GetAllClases obtiene todas las clases de una empresa
func (s *ClaseService) GetAllClases(empresaID uint64) ([]Clase, error) {
	return s.Repo.GetAllClases(empresaID)
}

// GetClase obtiene una clase específica
func (s *ClaseService) GetClase(id uint64, empresaID uint64) (Clase, error) {
	return s.Repo.GetClase(id, empresaID)
}

// CreateClase crea una nueva clase
func (s *ClaseService) CreateClase(clase *Clase) error {
	return s.Repo.CreateClase(clase)
}

// UpdateClase actualiza una clase existente y reindexa productos relacionados en Redis
func (s *ClaseService) UpdateClase(id uint64, empresaID uint64, clase *Clase) (Clase, error) {
	updated, err := s.Repo.UpdateClase(id, empresaID, clase)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		_ = s.ProductoSearch.UpdateProductsByClassification("clase", clase.Nombre)
	}
	return updated, nil
}

// DeleteClase elimina una clase
func (s *ClaseService) DeleteClase(id uint64, empresaID uint64) error {
	return s.Repo.DeleteClase(id, empresaID)
}

// PatchClase actualiza parcialmente una clase
func (s *ClaseService) PatchClase(id uint64, empresaID uint64, fields map[string]interface{}) (Clase, error) {
	updated, err := s.Repo.PatchClase(id, empresaID, fields)
	if err != nil {
		return updated, err
	}
	if s.ProductoSearch != nil {
		if nombre, ok := fields["nombre"].(string); ok {
			_ = s.ProductoSearch.UpdateProductsByClassification("clase", nombre)
		}
	}
	return updated, nil
}
