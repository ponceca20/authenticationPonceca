package users

import (
	"time"

	"gorm.io/gorm"
)

// LoginAttempt registra intentos de inicio de sesión para prevenir ataques de fuerza bruta
type LoginAttempt struct {
	gorm.Model
	Identifier string `gorm:"type:varchar(50);index"`
	IP         string `gorm:"type:varchar(45);index"`
	Success    bool   `gorm:"default:false;index"`
}

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
	// Relaciones
	Usuario *Usuario `json:"usuario,omitempty" gorm:"foreignKey:PersonaID"`
}

// Usuario representa la tabla USUARIO.
type Usuario struct {
	gorm.Model
	PersonaID      uint64   `json:"persona_id" gorm:"index"`
	Persona        *Persona `json:"persona,omitempty" gorm:"foreignKey:PersonaID"`
	PasswordHash   string   `json:"password_hash"`
	Activo         bool     `json:"activo"`
	CreadoPor      uint64   `json:"creado_por"`
	ActualizadoPor uint64   `json:"actualizado_por"`
	// Relaciones
	Sesiones        []*Sesion         `json:"sesiones,omitempty" gorm:"foreignKey:UsuarioID"`
	UsuarioEmpresas []*UsuarioEmpresa `json:"usuario_empresas,omitempty" gorm:"foreignKey:UsuarioID"`
}

// UsuarioEmpresa representa la tabla USUARIO_EMPRESA.
type UsuarioEmpresa struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id" gorm:"index"`
	Usuario         *Usuario  `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
	EmpresaID       uint64    `json:"empresa_id" gorm:"index"`
	RolID           uint64    `json:"rol_id" gorm:"index"`
	Rol             *Rol      `json:"rol,omitempty" gorm:"foreignKey:RolID"`
	FechaAsignacion time.Time `json:"fecha_asignacion"`
}

// Rol representa la tabla ROL.
type Rol struct {
	gorm.Model
	EmpresaID uint64 `json:"empresa_id" gorm:"index"`
	Codigo    string `json:"codigo"`
	Nombre    string `json:"nombre"`
	Activo    bool   `json:"activo"`
	// Relaciones
	UsuarioEmpresas []*UsuarioEmpresa `json:"usuario_empresas,omitempty" gorm:"foreignKey:RolID"`
	RolModulos      []*RolModulo      `json:"rol_modulos,omitempty" gorm:"foreignKey:RolID"`
}

// Sesion representa la tabla SESION.
type Sesion struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id" gorm:"index"`
	Usuario         *Usuario  `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
	Token           string    `json:"token" gorm:"type:text"`
	RefreshToken    string    `json:"refresh_token" gorm:"type:text"`
	FechaExpiracion time.Time `json:"fecha_expiracion"`
	Activa          bool      `json:"activa" gorm:"default:true"`
	IP              string    `json:"ip" gorm:"type:varchar(45)"`
}

// Modulo representa la tabla MODULO.
type Modulo struct {
	gorm.Model
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
	Ruta   string `json:"ruta"`
	Activo bool   `json:"activo"`
	Orden  int    `json:"orden"`
	// Relaciones
	RolModulos []*RolModulo `json:"rol_modulos,omitempty" gorm:"foreignKey:ModuloID"`
}

// RolModulo representa la tabla ROL_MODULO.
type RolModulo struct {
	gorm.Model
	RolID    uint64  `json:"rol_id" gorm:"index"`
	Rol      *Rol    `json:"rol,omitempty" gorm:"foreignKey:RolID"`
	ModuloID uint64  `json:"modulo_id" gorm:"index"`
	Modulo   *Modulo `json:"modulo,omitempty" gorm:"foreignKey:ModuloID"`
	Acceso   bool    `json:"acceso"`
}
