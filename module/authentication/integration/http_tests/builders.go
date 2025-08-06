package http_tests

import (
	"fmt"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/role"
	"time"
)

// =============================
// ORGANIZATION BUILDER
// =============================

// OrganizationBuilder ayuda a construir datos de organización de forma fluida
type OrganizationBuilder struct {
	suite *HTTPIntegrationTestSuite
	data  organization.OrganizationRegistrationDTO
}

// NewOrganizationBuilder crea un nuevo builder para organizaciones
func (s *HTTPIntegrationTestSuite) NewOrganizationBuilder(orgType string) *OrganizationBuilder {
	return &OrganizationBuilder{
		suite: s,
		data:  s.createOrganizationRegistrationData(orgType),
	}
}

// WithName establece el nombre de la organización
func (b *OrganizationBuilder) WithName(name string) *OrganizationBuilder {
	b.data.Name = name
	return b
}

// WithDescription establece la descripción de la organización
func (b *OrganizationBuilder) WithDescription(description string) *OrganizationBuilder {
	b.data.Description = description
	return b
}

// WithWebsite establece el sitio web de la organización
func (b *OrganizationBuilder) WithWebsite(website string) *OrganizationBuilder {
	b.data.Website = website
	return b
}

// WithFounder establece los datos del fundador
func (b *OrganizationBuilder) WithFounder(email, firstName, lastName, password string) *OrganizationBuilder {
	b.data.Identity = auth.RegisterDTO{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
	}
	return b
}

// WithFounderEmail establece solo el email del fundador (mantiene otros datos por defecto)
func (b *OrganizationBuilder) WithFounderEmail(email string) *OrganizationBuilder {
	b.data.Identity.Email = email
	return b
}

// Build construye y retorna los datos de la organización
func (b *OrganizationBuilder) Build() organization.OrganizationRegistrationDTO {
	return b.data
}

// =============================
// USER BUILDER
// =============================

// UserBuilder ayuda a construir datos de usuario de forma fluida
type UserBuilder struct {
	suite *HTTPIntegrationTestSuite
	data  auth.RegisterDTO
}

// NewUserBuilder crea un nuevo builder para usuarios
func (s *HTTPIntegrationTestSuite) NewUserBuilder() *UserBuilder {
	uniqueEmail := s.generateUniqueEmail("user")

	return &UserBuilder{
		suite: s,
		data: auth.RegisterDTO{
			Email:     uniqueEmail,
			FirstName: "Test",
			LastName:  "User",
			Password:  "TestPassword123!",
		},
	}
}

// WithEmail establece el email del usuario
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.data.Email = email
	return b
}

// WithName establece el nombre completo del usuario
func (b *UserBuilder) WithName(firstName, lastName string) *UserBuilder {
	b.data.FirstName = firstName
	b.data.LastName = lastName
	return b
}

// WithPassword establece la contraseña del usuario
func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.data.Password = password
	return b
}

// Build construye y retorna los datos del usuario
func (b *UserBuilder) Build() auth.RegisterDTO {
	return b.data
}

// BuildLoginDTO construye y retorna datos de login
func (b *UserBuilder) BuildLoginDTO() auth.LoginDTO {
	return auth.LoginDTO{
		Email:    b.data.Email,
		Password: b.data.Password,
	}
}

// =============================
// ROLE BUILDER
// =============================

// RoleBuilder ayuda a construir datos de rol de forma fluida
type RoleBuilder struct {
	suite *HTTPIntegrationTestSuite
	data  role.RoleDTO
}

// NewRoleBuilder crea un nuevo builder para roles
func (s *HTTPIntegrationTestSuite) NewRoleBuilder(name string) *RoleBuilder {
	return &RoleBuilder{
		suite: s,
		data: role.RoleDTO{
			Name:           name,
			DisplayName:    name,
			Description:    fmt.Sprintf("Test role: %s", name),
			HierarchyLevel: 50,
			Permissions:    []role.PermissionDTO{},
		},
	}
}

// WithDisplayName establece el nombre para mostrar del rol
func (b *RoleBuilder) WithDisplayName(displayName string) *RoleBuilder {
	b.data.DisplayName = displayName
	return b
}

// WithDescription establece la descripción del rol
func (b *RoleBuilder) WithDescription(description string) *RoleBuilder {
	b.data.Description = description
	return b
}

// WithHierarchyLevel establece el nivel jerárquico del rol
func (b *RoleBuilder) WithHierarchyLevel(level int) *RoleBuilder {
	b.data.HierarchyLevel = level
	return b
}

// WithPermission agrega un permiso al rol
func (b *RoleBuilder) WithPermission(resource string, actions []string, scope string) *RoleBuilder {
	permission := role.PermissionDTO{
		Resource: resource,
		Actions:  actions,
		Scope:    scope,
	}
	b.data.Permissions = append(b.data.Permissions, permission)
	return b
}

// WithReadPermission agrega permiso de lectura para un recurso
func (b *RoleBuilder) WithReadPermission(resource, scope string) *RoleBuilder {
	return b.WithPermission(resource, []string{"read"}, scope)
}

// WithWritePermission agrega permisos de lectura y escritura para un recurso
func (b *RoleBuilder) WithWritePermission(resource, scope string) *RoleBuilder {
	return b.WithPermission(resource, []string{"read", "create", "update"}, scope)
}

// WithFullPermission agrega todos los permisos para un recurso
func (b *RoleBuilder) WithFullPermission(resource, scope string) *RoleBuilder {
	return b.WithPermission(resource, []string{"read", "create", "update", "delete"}, scope)
}

// Build construye y retorna los datos del rol
func (b *RoleBuilder) Build() role.RoleDTO {
	return b.data
}

// =============================
// INVITATION BUILDER
// =============================

// InvitationBuilder ayuda a construir datos de invitación de forma fluida
type InvitationBuilder struct {
	suite *HTTPIntegrationTestSuite
	data  invitation.InvitationDTO
}

// NewInvitationBuilder crea un nuevo builder para invitaciones
func (s *HTTPIntegrationTestSuite) NewInvitationBuilder(roleID string) *InvitationBuilder {
	uniqueEmail := s.generateUniqueEmail("invite")

	return &InvitationBuilder{
		suite: s,
		data: invitation.InvitationDTO{
			Email:  uniqueEmail,
			RoleID: roleID,
		},
	}
}

// WithEmail establece el email de la invitación
func (b *InvitationBuilder) WithEmail(email string) *InvitationBuilder {
	b.data.Email = email
	return b
}

// WithRoleID establece el ID del rol para la invitación
func (b *InvitationBuilder) WithRoleID(roleID string) *InvitationBuilder {
	b.data.RoleID = roleID
	return b
}

// Build construye y retorna los datos de la invitación
func (b *InvitationBuilder) Build() invitation.InvitationDTO {
	return b.data
}

// =============================
// ACCEPT INVITATION BUILDER
// =============================

// AcceptInvitationBuilder ayuda a construir datos para aceptar invitaciones
type AcceptInvitationBuilder struct {
	suite *HTTPIntegrationTestSuite
	data  invitation.AcceptInvitationDTO
}

// NewAcceptInvitationBuilder crea un nuevo builder para aceptar invitaciones
func (s *HTTPIntegrationTestSuite) NewAcceptInvitationBuilder(token string) *AcceptInvitationBuilder {
	return &AcceptInvitationBuilder{
		suite: s,
		data: invitation.AcceptInvitationDTO{
			Token:     token,
			FirstName: "Test",
			LastName:  "Employee",
			Password:  "Employee2025!",
		},
	}
}

// WithName establece el nombre completo
func (b *AcceptInvitationBuilder) WithName(firstName, lastName string) *AcceptInvitationBuilder {
	b.data.FirstName = firstName
	b.data.LastName = lastName
	return b
}

// WithPassword establece la contraseña
func (b *AcceptInvitationBuilder) WithPassword(password string) *AcceptInvitationBuilder {
	b.data.Password = password
	return b
}

// Build construye y retorna los datos para aceptar invitación
func (b *AcceptInvitationBuilder) Build() invitation.AcceptInvitationDTO {
	return b.data
}

// =============================
// PREDEFINED ROLES
// =============================

// CorporateRoles contiene definiciones de roles corporativos comunes
type CorporateRoles struct {
	suite *HTTPIntegrationTestSuite
}

// NewCorporateRoles crea un generador de roles corporativos
func (s *HTTPIntegrationTestSuite) NewCorporateRoles() *CorporateRoles {
	return &CorporateRoles{suite: s}
}

// Manager crea un rol de gerente con permisos departamentales
func (cr *CorporateRoles) Manager() role.RoleDTO {
	return cr.suite.NewRoleBuilder("department_manager").
		WithDisplayName("Gerente de Departamento").
		WithDescription("Gerente con acceso departamental").
		WithHierarchyLevel(90).
		WithPermission("employees", []string{"read", "update", "invite"}, "department").
		WithPermission("reports", []string{"read", "create"}, "department").
		Build()
}

// SeniorDeveloper crea un rol de desarrollador senior
func (cr *CorporateRoles) SeniorDeveloper() role.RoleDTO {
	return cr.suite.NewRoleBuilder("senior_developer").
		WithDisplayName("Desarrollador Senior").
		WithDescription("Desarrollador con experiencia avanzada").
		WithHierarchyLevel(70).
		WithPermission("projects", []string{"read", "update"}, "department").
		WithPermission("code", []string{"read", "write", "review"}, "own").
		Build()
}

// Accountant crea un rol de contador
func (cr *CorporateRoles) Accountant() role.RoleDTO {
	return cr.suite.NewRoleBuilder("accountant").
		WithDisplayName("Contador").
		WithDescription("Responsable de la contabilidad").
		WithHierarchyLevel(60).
		WithPermission("invoices", []string{"read", "create", "update"}, "department").
		WithPermission("finances", []string{"read"}, "organization").
		Build()
}

// HRManager crea un rol de gerente de recursos humanos
func (cr *CorporateRoles) HRManager() role.RoleDTO {
	return cr.suite.NewRoleBuilder("hr_manager").
		WithDisplayName("Gerente de Recursos Humanos").
		WithDescription("Gestión de personal y recursos humanos").
		WithHierarchyLevel(85).
		WithPermission("employees", []string{"read", "create", "update", "invite"}, "organization").
		WithPermission("payroll", []string{"read", "create", "update"}, "organization").
		WithPermission("benefits", []string{"read", "create", "update"}, "organization").
		Build()
}

// SalesRepresentative crea un rol de representante de ventas
func (cr *CorporateRoles) SalesRepresentative() role.RoleDTO {
	return cr.suite.NewRoleBuilder("sales_representative").
		WithDisplayName("Representante de Ventas").
		WithDescription("Ejecutivo de ventas con acceso a clientes").
		WithHierarchyLevel(40).
		WithPermission("customers", []string{"read", "create", "update"}, "own").
		WithPermission("orders", []string{"read", "create", "update"}, "own").
		WithPermission("reports", []string{"read"}, "own").
		Build()
}

// =============================
// EMPLOYEE GENERATOR
// =============================

// Employee representa un empleado con sus datos completos
type Employee struct {
	Email     string
	RoleName  string
	Name      string
	FirstName string
	LastName  string
	Password  string
	RoleID    string
}

// EmployeeGenerator ayuda a generar conjuntos de empleados para testing
type EmployeeGenerator struct {
	suite *HTTPIntegrationTestSuite
}

// NewEmployeeGenerator crea un generador de empleados
func (s *HTTPIntegrationTestSuite) NewEmployeeGenerator() *EmployeeGenerator {
	return &EmployeeGenerator{suite: s}
}

// CreateCorporateTeam crea un equipo corporativo típico
func (eg *EmployeeGenerator) CreateCorporateTeam(createdRoles map[string]string) []Employee {
	employees := []Employee{
		{
			Email:     eg.suite.generateUniqueEmail("gerente.it"),
			RoleName:  "department_manager",
			Name:      "David Torres",
			FirstName: "David",
			LastName:  "Torres",
			Password:  "Employee2025!",
		},
		{
			Email:     eg.suite.generateUniqueEmail("dev.senior"),
			RoleName:  "senior_developer",
			Name:      "Miguel Herrera",
			FirstName: "Miguel",
			LastName:  "Herrera",
			Password:  "Employee2025!",
		},
		{
			Email:     eg.suite.generateUniqueEmail("contador"),
			RoleName:  "accountant",
			Name:      "Ana García",
			FirstName: "Ana",
			LastName:  "García",
			Password:  "Employee2025!",
		},
	}

	// Asignar role IDs
	for i := range employees {
		if roleID, exists := createdRoles[employees[i].RoleName]; exists {
			employees[i].RoleID = roleID
		}
	}

	return employees
}

// CreateSalesTeam crea un equipo de ventas
func (eg *EmployeeGenerator) CreateSalesTeam(createdRoles map[string]string) []Employee {
	employees := []Employee{
		{
			Email:     eg.suite.generateUniqueEmail("sales.lead"),
			RoleName:  "sales_manager",
			Name:      "Carlos Vendedor",
			FirstName: "Carlos",
			LastName:  "Vendedor",
			Password:  "Sales2025!",
		},
		{
			Email:     eg.suite.generateUniqueEmail("sales.rep1"),
			RoleName:  "sales_representative",
			Name:      "María Comercial",
			FirstName: "María",
			LastName:  "Comercial",
			Password:  "Sales2025!",
		},
		{
			Email:     eg.suite.generateUniqueEmail("sales.rep2"),
			RoleName:  "sales_representative",
			Name:      "Pedro Cliente",
			FirstName: "Pedro",
			LastName:  "Cliente",
			Password:  "Sales2025!",
		},
	}

	// Asignar role IDs
	for i := range employees {
		if roleID, exists := createdRoles[employees[i].RoleName]; exists {
			employees[i].RoleID = roleID
		}
	}

	return employees
}

// =============================
// TIMESTAMPS HELPER
// =============================

// TimeHelper ayuda con la generación de timestamps para testing
type TimeHelper struct{}

// NewTimeHelper crea un helper para timestamps
func (s *HTTPIntegrationTestSuite) NewTimeHelper() *TimeHelper {
	return &TimeHelper{}
}

// FutureTime retorna un tiempo futuro
func (th *TimeHelper) FutureTime(hours int) time.Time {
	return time.Now().Add(time.Duration(hours) * time.Hour)
}

// PastTime retorna un tiempo pasado
func (th *TimeHelper) PastTime(hours int) time.Time {
	return time.Now().Add(-time.Duration(hours) * time.Hour)
}

// Timestamp retorna timestamp actual en formato RFC3339
func (th *TimeHelper) Timestamp() string {
	return time.Now().Format(time.RFC3339)
}
