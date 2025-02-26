package users

import (
	"time"

	"gorm.io/gorm"
)

// Persona representa la tabla PERSONA.
type Persona struct {
	gorm.Model
	DocumentoTipo   string    `json:"documento_tipo"`
	DocumentoNumero string    `json:"documento_numero"`
	Foto            string    `json:"foto"`
	Nombre          string    `json:"nombre"`
	Apellidos       string    `json:"apellidos"`
	Email           string    `json:"email"`
	Telefono        string    `json:"telefono"`
	Direccion       string    `json:"direccion"`
	Ciudad          string    `json:"ciudad"`
	Pais            string    `json:"pais"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
}

// Usuario representa la tabla USUARIO.
type Usuario struct {
	gorm.Model
	PersonaID      uint64 `json:"persona_id"`
	PasswordHash   string `json:"password_hash"`
	Activo         bool   `json:"activo"`
	CreadoPor      uint64 `json:"creado_por"`
	ActualizadoPor uint64 `json:"actualizado_por"`
}

// UsuarioEmpresa representa la tabla USUARIO_EMPRESA.
type UsuarioEmpresa struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id"`
	EmpresaID       uint64    `json:"empresa_id"`
	RolID           uint64    `json:"rol_id"`
	FechaAsignacion time.Time `json:"fecha_asignacion"`
}

// Rol representa la tabla ROL.
type Rol struct {
	gorm.Model
	EmpresaID uint64 `json:"empresa_id"`
	Codigo    string `json:"codigo"`
	Nombre    string `json:"nombre"`
	Activo    bool   `json:"activo"`
}

// Sesion representa la tabla SESION.
// Se eliminó FechaCreacion ya que gorm.Model provee CreatedAt.
type Sesion struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id"`
	Token           string    `json:"token"`
	RefreshToken    string    `json:"refresh_token"`
	FechaExpiracion time.Time `json:"fecha_expiracion"`
	Activa          bool      `json:"activa"`
}

// Modulo representa la tabla MODULO.
type Modulo struct {
	gorm.Model
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
	Ruta   string `json:"ruta"`
	Activo bool   `json:"activo"`
	Orden  int    `json:"orden"`
}

// RolModulo representa la tabla ROL_MODULO.
// Se eliminaron FechaCreacion y FechaActualizacion ya que gorm.Model provee CreatedAt y UpdatedAt.
type RolModulo struct {
	gorm.Model
	RolID    uint64 `json:"rol_id"`
	ModuloID uint64 `json:"modulo_id"`
	Acceso   bool   `json:"acceso"`
}
