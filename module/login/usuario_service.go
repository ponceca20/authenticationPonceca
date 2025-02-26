package users

import "golang.org/x/crypto/bcrypt"

// UsuarioService contiene la lógica de negocio para usuarios
type UsuarioService struct {
	Repo UsuarioRepository
}

// NewUsuarioService crea una nueva instancia del servicio
func NewUsuarioService(repo UsuarioRepository) *UsuarioService {
	return &UsuarioService{Repo: repo}
}

// GetAllUsuarios obtiene todos los usuarios
func (s *UsuarioService) GetAllUsuarios() ([]Usuario, error) {
	return s.Repo.GetAllUsuarios()
}

// GetUsuario obtiene un usuario por ID delegando en el repositorio
func (s *UsuarioService) GetUsuario(id uint64) (Usuario, error) {
	return s.Repo.GetUsuario(id)
}

// hashPassword genera el hash para una contraseña en texto plano.
func (s *UsuarioService) hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CreateUsuario crea un nuevo usuario
func (s *UsuarioService) CreateUsuario(usuario *Usuario) error {
	// Hash de la contraseña si se proporcionó
	if usuario.PasswordHash != "" {
		hashedPassword, err := s.hashPassword(usuario.PasswordHash)
		if err != nil {
			return err
		}
		usuario.PasswordHash = hashedPassword
	}
	return s.Repo.CreateUsuario(usuario)
}

// UpdateUsuario actualiza un usuario existente
func (s *UsuarioService) UpdateUsuario(id uint64, usuario *Usuario) (Usuario, error) {
	// Si se actualiza la contraseña, generar hash
	if usuario.PasswordHash != "" {
		hashedPassword, err := s.hashPassword(usuario.PasswordHash)
		if err != nil {
			return Usuario{}, err
		}
		usuario.PasswordHash = hashedPassword
	}
	return s.Repo.UpdateUsuario(id, usuario)
}

// DeleteUsuario elimina un usuario
func (s *UsuarioService) DeleteUsuario(id uint64) error {
	return s.Repo.DeleteUsuario(id)
}

// SeedUsuarios inicializa múltiples usuarios
func (s *UsuarioService) SeedUsuarios(usuarios []Usuario) error {
	// TODO: Validar duplicados antes de insertar
	return s.Repo.SeedUsuarios(usuarios)
}

// PatchUsuario actualiza parcialmente un usuario
func (s *UsuarioService) PatchUsuario(id uint64, fields map[string]interface{}) (Usuario, error) {
	if password, ok := fields["password"].(string); ok && password != "" {
		hashedPassword, err := s.hashPassword(password)
		if err != nil {
			return Usuario{}, err
		}
		fields["password_hash"] = hashedPassword
		delete(fields, "password")
	} else if password, ok := fields["password_hash"].(string); ok && password != "" {
		hashedPassword, err := s.hashPassword(password)
		if err != nil {
			return Usuario{}, err
		}
		fields["password_hash"] = hashedPassword
	}
	return s.Repo.PatchUsuario(id, fields)
}
