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

// TipoDocumentoSunat representa la tabla TIPO_DOCUMENTO_SUNAT
type TipoDocumentoSunat struct {
	gorm.Model
	Codigo string `gorm:"type:varchar(10);unique;not null"`
	Nombre string `gorm:"type:varchar(100);not null"`
	// Se reintroduce la relación para que GORM genere la FK
	Personas []*Persona `gorm:"foreignKey:TipoDocumentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}

// Persona representa la tabla PERSONA
type Persona struct {
	gorm.Model
	TipoDocumentoID    uint64    `gorm:"not null"`
	DocumentoNumero    string    `gorm:"type:varchar(20)"`
	Foto               string    `gorm:"type:varchar(255)"`
	Nombre             string    `gorm:"type:varchar(100)"`
	Apellidos          string    `gorm:"type:varchar(100)"`
	Email              string    `gorm:"type:varchar(100);uniqueIndex"`
	Telefono           string    `gorm:"type:varchar(20)"`
	TelefonoSecundario string    `gorm:"type:varchar(20)"`
	Direccion          string    `gorm:"type:varchar(255)"`
	FechaNacimiento    time.Time `gorm:"type:date"`
	// Se reintroduce la relación para crear la restricción FK
	TipoDocumento TipoDocumentoSunat `gorm:"foreignKey:TipoDocumentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	// Relación uno a uno con Usuario
	Usuario *Usuario `gorm:"foreignKey:PersonaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

// Usuario representa la tabla USUARIO
type Usuario struct {
	gorm.Model
	PersonaID      uint64 `json:"persona_id" gorm:"uniqueIndex"`
	PasswordHash   string `json:"password_hash" gorm:"type:varchar(255);not null"`
	Activo         bool   `json:"activo" gorm:"default:true"`
	CreadoPor      uint64 `json:"creado_por"`
	ActualizadoPor uint64 `json:"actualizado_por"`
	// Relación con Persona (un usuario tiene una persona)
	Persona Persona `gorm:"foreignKey:PersonaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	// Un usuario puede tener muchas sesiones y relaciones con empresas
	Sesiones        []*Sesion         `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UsuarioEmpresas []*UsuarioEmpresa `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Empresa representa la tabla EMPRESA (implícita en la relación USUARIO_EMPRESA y ROL)
type Empresa struct {
	gorm.Model
	// Campos básicos de empresa
	Nombre    string `json:"nombre" gorm:"type:varchar(200);not null"`
	RUC       string `json:"ruc" gorm:"type:varchar(11);uniqueIndex"`
	Direccion string `json:"direccion" gorm:"type:varchar(255)"`
	Activo    bool   `json:"activo" gorm:"default:true"`
	// Relaciones
	UsuarioEmpresas []*UsuarioEmpresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Roles           []*Rol            `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// UsuarioEmpresa representa la tabla USUARIO_EMPRESA
type UsuarioEmpresa struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id" gorm:"index;not null"`
	EmpresaID       uint64    `json:"empresa_id" gorm:"index;not null"`
	RolID           uint64    `json:"rol_id" gorm:"index;not null"`
	FechaAsignacion time.Time `json:"fecha_asignacion" gorm:"default:CURRENT_TIMESTAMP"`
	// Relación con Usuario
	Usuario Usuario `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	// Relación con Empresa
	Empresa Empresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	// Relación con Rol
	Rol Rol `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
}

// Rol representa la tabla ROL
type Rol struct {
	gorm.Model
	EmpresaID uint64 `json:"empresa_id" gorm:"index;not null"`
	Codigo    string `json:"codigo" gorm:"type:varchar(50);not null"`
	Nombre    string `json:"nombre" gorm:"type:varchar(100);not null"`
	Activo    bool   `json:"activo" gorm:"default:true"`
	// Relación con Empresa
	Empresa Empresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	// Un rol puede estar relacionado a muchos UsuarioEmpresa y RolModulo
	UsuarioEmpresas []*UsuarioEmpresa `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	RolModulos      []*RolModulo      `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Sesion representa la tabla SESION
type Sesion struct {
	gorm.Model
	UsuarioID       uint64    `json:"usuario_id" gorm:"index;not null"`
	Token           string    `json:"token" gorm:"type:text;not null"`
	RefreshToken    string    `json:"refresh_token" gorm:"type:text"`
	FechaExpiracion time.Time `json:"fecha_expiracion" gorm:"not null"`
	Activa          bool      `json:"activa" gorm:"default:true"`
	IP              string    `json:"ip" gorm:"type:varchar(45)"`
	// Relación con Usuario
	Usuario Usuario `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Modulo representa la tabla MODULO
type Modulo struct {
	gorm.Model
	Codigo string `json:"codigo" gorm:"type:varchar(50);unique;not null"`
	Nombre string `json:"nombre" gorm:"type:varchar(100);not null"`
	Ruta   string `json:"ruta" gorm:"type:varchar(255)"`
	Activo bool   `json:"activo" gorm:"default:true"`
	Orden  int    `json:"orden" gorm:"default:0"`
	// Un módulo puede estar relacionado a muchos RolModulo
	RolModulos []*RolModulo `gorm:"foreignKey:ModuloID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// RolModulo representa la tabla ROL_MODULO
type RolModulo struct {
	gorm.Model
	RolID    uint64 `json:"rol_id" gorm:"index;not null"`
	ModuloID uint64 `json:"modulo_id" gorm:"index;not null"`
	Acceso   bool   `json:"acceso" gorm:"default:false"`
	// Relaciones para materializar la FK
	Rol    Rol    `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Modulo Modulo `gorm:"foreignKey:ModuloID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
