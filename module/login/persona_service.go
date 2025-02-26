package users

// PersonaService contiene la lógica de negocio para Persona.
type PersonaService struct {
	Repo PersonaRepository
}

func NewPersonaService(repo PersonaRepository) *PersonaService {
	return &PersonaService{Repo: repo}
}

func (s *PersonaService) GetAllPersonas() ([]Persona, error) {
	return s.Repo.GetAllPersonas()
}

func (s *PersonaService) GetPersona(id uint64) (Persona, error) {
	return s.Repo.GetPersona(id)
}

func (s *PersonaService) CreatePersona(persona *Persona) error {
	return s.Repo.CreatePersona(persona)
}

func (s *PersonaService) UpdatePersona(id uint64, persona *Persona) (Persona, error) {
	return s.Repo.UpdatePersona(id, persona)
}

func (s *PersonaService) DeletePersona(id uint64) error {
	return s.Repo.DeletePersona(id)
}

func (s *PersonaService) SeedPersonas(personas []Persona) error {
	return s.Repo.SeedPersonas(personas)
}

func (s *PersonaService) PatchPersona(id uint64, fields map[string]interface{}) (Persona, error) {
	return s.Repo.PatchPersona(id, fields)
}
