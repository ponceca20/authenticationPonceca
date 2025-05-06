package product

// ComentarioProductoService contiene la lógica de negocio para ComentarioProducto
type ComentarioProductoService struct {
	Repo ComentarioProductoRepository
}

// NewComentarioProductoService crea una nueva instancia del servicio
func NewComentarioProductoService(repo ComentarioProductoRepository) *ComentarioProductoService {
	return &ComentarioProductoService{Repo: repo}
}

// GetAllComentariosProducto obtiene todos los comentarios para un producto específico
func (s *ComentarioProductoService) GetAllComentariosProducto(productoID *uint64) ([]ComentarioProducto, error) {
	return s.Repo.GetAllComentariosProducto(productoID)
}

// GetComentariosGenericos obtiene todos los comentarios genéricos
func (s *ComentarioProductoService) GetComentariosGenericos() ([]ComentarioProducto, error) {
	return s.Repo.GetComentariosGenericos()
}

// GetComentarioProducto obtiene un comentario específico
func (s *ComentarioProductoService) GetComentarioProducto(id uint64) (ComentarioProducto, error) {
	return s.Repo.GetComentarioProducto(id)
}

// CreateComentarioProducto crea un nuevo comentario
func (s *ComentarioProductoService) CreateComentarioProducto(comentario *ComentarioProducto) error {
	return s.Repo.CreateComentarioProducto(comentario)
}

// UpdateComentarioProducto actualiza un comentario existente
func (s *ComentarioProductoService) UpdateComentarioProducto(id uint64, comentario *ComentarioProducto) (ComentarioProducto, error) {
	return s.Repo.UpdateComentarioProducto(id, comentario)
}

// DeleteComentarioProducto elimina un comentario
func (s *ComentarioProductoService) DeleteComentarioProducto(id uint64) error {
	return s.Repo.DeleteComentarioProducto(id)
}

// SeedComentariosProducto inicializa múltiples comentarios
func (s *ComentarioProductoService) SeedComentariosProducto(comentarios []ComentarioProducto) error {
	return s.Repo.SeedComentariosProducto(comentarios)
}

// PatchComentarioProducto actualiza parcialmente un comentario
func (s *ComentarioProductoService) PatchComentarioProducto(id uint64, fields map[string]interface{}) (ComentarioProducto, error) {
	return s.Repo.PatchComentarioProducto(id, fields)
}
