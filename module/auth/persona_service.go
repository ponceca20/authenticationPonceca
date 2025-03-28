package auth

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

func (s *PersonaService) GetPersona(id uint) (Persona, error) { // Changed from uint64 to uint
	return s.Repo.GetPersona(id)
}

func (s *PersonaService) CreatePersona(persona *Persona) error {
	return s.Repo.CreatePersona(persona)
}

func (s *PersonaService) UpdatePersona(id uint, persona *Persona) (Persona, error) { // Changed from uint64 to uint
	return s.Repo.UpdatePersona(id, persona)
}

func (s *PersonaService) DeletePersona(id uint) error { // Changed from uint64 to uint
	return s.Repo.DeletePersona(id)
}

func (s *PersonaService) SeedPersonas(personas []Persona) error {
	return s.Repo.SeedPersonas(personas)
}

func (s *PersonaService) PatchPersona(id uint, fields map[string]interface{}) (Persona, error) { // Changed from uint64 to uint
	return s.Repo.PatchPersona(id, fields)
}
