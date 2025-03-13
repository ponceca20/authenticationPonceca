package empresa

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Empresa representa la tabla emp_EMPRESA.
type Empresa struct {
	gorm.Model
	Codigo          string          `json:"codigo" gorm:"uniqueIndex:idx_empresa_codigo"`
	Ruc             string          `json:"ruc" gorm:"uniqueIndex:idx_empresa_ruc"`
	RazonSocial     string          `json:"razon_social"`
	NombreComercial string          `json:"nombre_comercial"`
	DireccionFiscal string          `json:"direccion_fiscal"`
	Telefono        string          `json:"telefono"`
	Email           string          `json:"email"`
	ZonaHoraria     string          `json:"zona_horaria"`
	DatosSunat      json.RawMessage `json:"datos_sunat"`
	DatosContacto   json.RawMessage `json:"datos_contacto"`
	Logos           json.RawMessage `json:"logos"`
	Latitud         float64         `json:"latitud"`
	Longitud        float64         `json:"longitud"`
	Horario         json.RawMessage `json:"horario"`
	Estado          string          `json:"estado" gorm:"default:'ACTIVO'"` // ACTIVO|INACTIVO|ELIMINADO

	// Relaciones (agregadas)
	Configuraciones    []Configuracion   `json:"configuraciones,omitempty" gorm:"foreignKey:EmpresaID"`
	RelacionesMatriz   []EmpresaRelacion `json:"relaciones_matriz,omitempty" gorm:"foreignKey:EmpresaMatrizID"`
	RelacionesSucursal []EmpresaRelacion `json:"relaciones_sucursal,omitempty" gorm:"foreignKey:EmpresaSucursalID"`
	Suscripciones      []SuscSuscripcion `json:"suscripciones,omitempty" gorm:"foreignKey:EmpresaID"`
	Documentos         []SuscDocumento   `json:"documentos,omitempty" gorm:"foreignKey:EmpresaID"`
}

// EmpresaRelacion representa la tabla emp_EMPRESA_RELACION.
type EmpresaRelacion struct {
	gorm.Model
	EmpresaMatrizID   uint64  `json:"empresa_matriz_id" gorm:"index"`
	EmpresaMatriz     Empresa `json:"-" gorm:"foreignKey:EmpresaMatrizID"`
	EmpresaSucursalID uint64  `json:"empresa_sucursal_id" gorm:"index"`
	EmpresaSucursal   Empresa `json:"-" gorm:"foreignKey:EmpresaSucursalID"`
	TipoRelacion      string  `json:"tipo_relacion"` // SUCURSAL|FRANQUICIA|ASOCIADA
	Estado            string  `json:"estado" gorm:"default:'ACTIVO'"`
}

// Configuracion representa la tabla emp_CONFIGURACION.
type Configuracion struct {
	gorm.Model
	EmpresaID uint64  `json:"empresa_id" gorm:"index"`
	Empresa   Empresa `json:"-" gorm:"foreignKey:EmpresaID"`
	Clave     string  `json:"clave" gorm:"uniqueIndex:idx_empresa_config"`
	Valor     string  `json:"valor"`
	EsSecreto bool    `json:"es_secreto" gorm:"default:false"`
}

// SuscPlan representa la tabla emp_SUSC_PLAN.
type SuscPlan struct {
	gorm.Model
	Nombre        string  `json:"nombre"`
	Descripcion   string  `json:"descripcion"`
	PrecioMensual float64 `json:"precio_mensual" gorm:"type:decimal(16,4)"`
	PrecioAnual   float64 `json:"precio_anual" gorm:"type:decimal(16,4)"`
	EsFree        bool    `json:"es_free" gorm:"default:false"`
	Activo        bool    `json:"activo" gorm:"default:true"`

	// Relación (agregada)
	Suscripciones []SuscSuscripcion `json:"suscripciones,omitempty" gorm:"foreignKey:PlanID"`
}

// SuscSuscripcion representa la tabla emp_SUSC_SUSCRIPCION.
type SuscSuscripcion struct {
	gorm.Model
	EmpresaID      uint64    `json:"empresa_id" gorm:"index"`
	Empresa        Empresa   `json:"-" gorm:"foreignKey:EmpresaID"`
	PlanID         uint64    `json:"plan_id" gorm:"index"`
	Plan           SuscPlan  `json:"-" gorm:"foreignKey:PlanID"`
	FechaInicio    time.Time `json:"fecha_inicio"`
	FechaFin       time.Time `json:"fecha_fin"`
	Periodicidad   string    `json:"periodicidad"` // FREE|MENSUAL|ANUAL
	Estado         string    `json:"estado"`       // ACTIVA|PENDIENTE|SUSPENDIDA|CANCELADA
	Notificado     bool      `json:"notificado" gorm:"default:false"`
	LimiteUsuarios int       `json:"limite_usuarios"`

	// Relaciones (agregadas)
	Pagos      []SuscPago      `json:"pagos,omitempty" gorm:"foreignKey:SuscripcionID"`
	Documentos []SuscDocumento `json:"documentos,omitempty" gorm:"foreignKey:SuscripcionID"`
}

// SuscPago representa la tabla emp_SUSC_PAGO.
type SuscPago struct {
	gorm.Model
	SuscripcionID  uint64          `json:"suscripcion_id" gorm:"index"`
	Suscripcion    SuscSuscripcion `json:"-" gorm:"foreignKey:SuscripcionID"`
	Monto          float64         `json:"monto" gorm:"type:decimal(16,4)"`
	FechaPago      time.Time       `json:"fecha_pago"`
	PeriodoInicio  time.Time       `json:"periodo_inicio"`
	PeriodoFin     time.Time       `json:"periodo_fin"`
	Estado         string          `json:"estado"`      // PENDIENTE|PAGADO|CANCELADO
	MetodoPago     string          `json:"metodo_pago"` // EFECTIVO|TARJETA|TRANSFERENCIA|DEPOSITO
	ReferenciaPago string          `json:"referencia_pago"`
}

// SuscDocumento representa la tabla emp_SUSC_DOCUMENTO.
type SuscDocumento struct {
	gorm.Model
	EmpresaID        uint64          `json:"empresa_id" gorm:"index"`
	Empresa          Empresa         `json:"-" gorm:"foreignKey:EmpresaID"`
	SuscripcionID    uint64          `json:"suscripcion_id" gorm:"index"`
	Suscripcion      SuscSuscripcion `json:"-" gorm:"foreignKey:SuscripcionID"`
	Tipo             string          `json:"tipo"` // FACTURA|BOLETA|NOTA_CREDITO|NOTA_DEBITO
	MontoTotal       float64         `json:"monto_total" gorm:"type:decimal(16,4)"`
	MontoPendiente   float64         `json:"monto_pendiente" gorm:"type:decimal(16,4)"`
	FechaEmision     time.Time       `json:"fecha_emision"`
	FechaVencimiento time.Time       `json:"fecha_vencimiento"`
	EsCredito        bool            `json:"es_credito" gorm:"default:false"`
	EstadoPago       string          `json:"estado_pago"` // PENDIENTE|PAGADO|VENCIDO|ANULADO
}
