package auth

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

// TipoDocumentoSunat representa la tabla TIPO_DOCUMENTO_SUNAT
type TipoDocumentoSunat struct {
	gorm.Model
	Codigo string `json:"codigo" gorm:"type:varchar(20);index"`
	Nombre string `json:"nombre" gorm:"type:varchar(100)"`
	// Relaciones
	Personas []*Persona `json:"personas,omitempty" gorm:"foreignKey:TipoDocumentoID"`
}

// Persona representa la tabla PERSONA.
type Persona struct {
	gorm.Model
	TipoDocumentoID    uint64              `json:"tipo_documento_id" gorm:"index"`
	TipoDocumento      *TipoDocumentoSunat `json:"tipo_documento,omitempty" gorm:"foreignKey:TipoDocumentoID"`
	DocumentoNumero    string              `json:"documento_numero" gorm:"type:varchar(20);index"`
	Foto               string              `json:"foto" gorm:"type:varchar(255)"`
	Nombre             string              `json:"nombre" gorm:"type:varchar(255)"`
	Apellidos          string              `json:"apellidos" gorm:"type:varchar(255)"`
	Email              string              `json:"email" gorm:"type:varchar(255);index"`
	Telefono           string              `json:"telefono" gorm:"type:varchar(20)"`
	TelefonoSecundario string              `json:"telefono_secundario" gorm:"type:varchar(20)"`
	Direccion          string              `json:"direccion" gorm:"type:text"`
	FechaNacimiento    time.Time           `json:"fecha_nacimiento" gorm:"type:date"`
	// Relaciones
	Usuario *Usuario `json:"usuario,omitempty" gorm:"foreignKey:PersonaID"`
}

// Usuario representa la tabla USUARIO.
type Usuario struct {
	gorm.Model
	PersonaID      uint64   `json:"persona_id" gorm:"index"`
	Persona        *Persona `json:"persona,omitempty" gorm:"foreignKey:PersonaID"`
	PasswordHash   string   `json:"password_hash" gorm:"type:varchar(250)"` // Para bcrypt
	Activo         bool     `json:"activo" gorm:"default:true;index"`
	CreadoPor      uint64   `json:"creado_por" gorm:"index"`
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
	ActiveSesion    bool      `json:"activo" gorm:"default:true;index"`
}

// Rol representa la tabla ROL.
type Rol struct {
	gorm.Model
	EmpresaID uint64 `json:"empresa_id" gorm:"index"`
	Codigo    string `json:"codigo" gorm:"type:varchar(50);uniqueIndex"`
	Nombre    string `json:"nombre" gorm:"type:varchar(250)"`
	Activo    bool   `json:"activo" gorm:"default:true;index"`
	// Relaciones
	UsuarioEmpresas []*UsuarioEmpresa `json:"usuario_empresas,omitempty" gorm:"foreignKey:RolID"`
	RolModulos      []*RolModulo      `json:"rol_modulos,omitempty" gorm:"foreignKey:RolID"`
}

// Sesion representa la tabla SESION.
type Sesion struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id" gorm:"index"`
	Usuario         *Usuario  `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
	Token           string    `json:"token" gorm:"type:varchar(500)"`
	RefreshToken    string    `json:"refresh_token" gorm:"type:varchar(500)"`
	FechaExpiracion time.Time `json:"fecha_expiracion" gorm:"type:timestamp;index"`
	Activa          bool      `json:"activa" gorm:"default:true;index"`
	IP              string    `json:"ip" gorm:"type:varchar(45);index"`
}

// Modulo representa la tabla MODULO.
type Modulo struct {
	gorm.Model
	Codigo string `json:"codigo" gorm:"type:varchar(50);uniqueIndex"`
	Nombre string `json:"nombre" gorm:"type:varchar(100)"`
	Ruta   string `json:"ruta" gorm:"type:varchar(255)"`
	Activo bool   `json:"activo" gorm:"default:true;index"`
	Orden  int    `json:"orden" gorm:"type:int"`
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
