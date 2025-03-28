package auth

// TipoDocumentoSunatService contiene la lógica de negocio para TipoDocumentoSunat.
type TipoDocumentoSunatService struct {
	Repo TipoDocumentoSunatRepository
}

func NewTipoDocumentoSunatService(repo TipoDocumentoSunatRepository) *TipoDocumentoSunatService {
	return &TipoDocumentoSunatService{Repo: repo}
}

func (s *TipoDocumentoSunatService) GetAllTiposDocumento() ([]TipoDocumentoSunat, error) {
	return s.Repo.GetAllTiposDocumento()
}

func (s *TipoDocumentoSunatService) GetTipoDocumento(id uint) (TipoDocumentoSunat, error) {
	return s.Repo.GetTipoDocumento(id)
}

func (s *TipoDocumentoSunatService) GetTipoDocumentoByCodigo(codigo string) (TipoDocumentoSunat, error) {
	return s.Repo.GetTipoDocumentoByCodigo(codigo)
}

func (s *TipoDocumentoSunatService) CreateTipoDocumento(tipoDocumento *TipoDocumentoSunat) error {
	return s.Repo.CreateTipoDocumento(tipoDocumento)
}

func (s *TipoDocumentoSunatService) UpdateTipoDocumento(id uint, tipoDocumento *TipoDocumentoSunat) (TipoDocumentoSunat, error) {
	return s.Repo.UpdateTipoDocumento(id, tipoDocumento)
}

func (s *TipoDocumentoSunatService) DeleteTipoDocumento(id uint) error {
	return s.Repo.DeleteTipoDocumento(id)
}

func (s *TipoDocumentoSunatService) SeedTiposDocumento(tiposDocumento []TipoDocumentoSunat) error {
	return s.Repo.SeedTiposDocumento(tiposDocumento)
}

func (s *TipoDocumentoSunatService) PatchTipoDocumento(id uint, fields map[string]interface{}) (TipoDocumentoSunat, error) {
	return s.Repo.PatchTipoDocumento(id, fields)
}
