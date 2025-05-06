package product

// CategoriaService contiene la lógica de negocio para Categoria
type CategoriaService struct {
    Repo            CategoriaRepository
    ProductoSearch  *ProductoSearchService
}

// NewCategoriaService crea una nueva instancia del servicio
func NewCategoriaService(repo CategoriaRepository, search *ProductoSearchService) *CategoriaService {
    return &CategoriaService{Repo: repo, ProductoSearch: search}
}

// GetAllCategorias obtiene todas las categorias de una empresa
func (s *CategoriaService) GetAllCategorias(empresaID uint64) ([]Categoria, error) {
    return s.Repo.GetAllCategorias(empresaID)
}

// GetCategoria obtiene una categoria específica
func (s *CategoriaService) GetCategoria(id uint64, empresaID uint64) (Categoria, error) {
    return s.Repo.GetCategoria(id, empresaID)
}

// CreateCategoria crea una nueva categoria
func (s *CategoriaService) CreateCategoria(categoria *Categoria) error {
    return s.Repo.CreateCategoria(categoria)
}

// UpdateCategoria actualiza una categoria existente
func (s *CategoriaService) UpdateCategoria(id uint64, empresaID uint64, categoria *Categoria) (Categoria, error) {
    updated, err := s.Repo.UpdateCategoria(id, empresaID, categoria)
    if err != nil {
        return updated, err
    }
    if s.ProductoSearch != nil {
        _ = s.ProductoSearch.UpdateProductsByClassification("categoria", categoria.Nombre)
    }
    return updated, nil
}

// DeleteCategoria elimina una categoria
func (s *CategoriaService) DeleteCategoria(id uint64, empresaID uint64) error {
    return s.Repo.DeleteCategoria(id, empresaID)
}

// PatchCategoria actualiza parcialmente una categoria
func (s *CategoriaService) PatchCategoria(id uint64, empresaID uint64, fields map[string]interface{}) (Categoria, error) {
    updated, err := s.Repo.PatchCategoria(id, empresaID, fields)
    if err != nil {
        return updated, err
    }
    if s.ProductoSearch != nil {
        if nombre, ok := fields["nombre"].(string); ok {
            _ = s.ProductoSearch.UpdateProductsByClassification("categoria", nombre)
        }
    }
    return updated, nil
}
