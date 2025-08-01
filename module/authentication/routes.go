package authentication

import (
	"practicev2/database"
	"practicev2/module/authentication/admin"
	"practicev2/module/authentication/audit"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/customer"
	"practicev2/module/authentication/department"
	"practicev2/module/authentication/guest"
	"practicev2/module/authentication/integration"
	"practicev2/module/authentication/invitation"
	"practicev2/module/authentication/middleware"
	"practicev2/module/authentication/organization"
	"practicev2/module/authentication/profile"
	"practicev2/module/authentication/role"
	"practicev2/module/authentication/unified"
	"practicev2/module/authentication/user"
	"practicev2/module/authentication/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all the routes for the authentication module.
func RegisterRoutes(app *fiber.App) {
	// --- DEPENDENCY INJECTION ---
	db := database.DBconn
	jwtService := utils.NewJWTService()

	// Repositories
	authRepo := auth.NewAuthRepository(db)
	orgRepo := organization.NewOrganizationRepository(db)
	userRepo := user.NewUserRepository(db)
	roleRepo := role.NewRoleRepository(db)
	profileRepo := profile.NewProfileRepository(db)
	invitationRepo := invitation.NewInvitationRepository(db)
	customerRepo := customer.NewCustomerRepository(db)
	guestRepo := guest.NewGuestRepository(db)
	departmentRepo := department.NewDepartmentRepository(db)
	auditRepo := audit.NewAuditRepository(db)
	unifiedRepo := unified.NewUnifiedRepository(db)
	_ = integration.NewIntegrationRepository(db) // Placeholder
	adminRepo := admin.NewAdminRepository(db)

	// Services
	authService := auth.NewAuthService(authRepo, jwtService)
	orgService := organization.NewOrganizationService(orgRepo, authRepo)
	userService := user.NewUserService(userRepo)
	roleService := role.NewRoleService(roleRepo)
	profileService := profile.NewProfileService(profileRepo, db)
	invitationService := invitation.NewInvitationService(invitationRepo)
	guestService := guest.NewGuestService(guestRepo)
	customerService := customer.NewCustomerService(customerRepo, authRepo, guestService)
	departmentService := department.NewDepartmentService(departmentRepo)
	auditService := audit.NewAuditService(auditRepo)
	unifiedService := unified.NewUnifiedService(unifiedRepo)
	integrationService := integration.NewIntegrationService(jwtService)
	adminService := admin.NewAdminService(adminRepo)

	// Handlers
	authHandler := auth.NewAuthHandler(authService)
	orgHandler := organization.NewOrganizationHandler(orgService)
	userHandler := user.NewUserHandler(userService)
	roleHandler := role.NewRoleHandler(roleService, userService) // Pass userService for role assignments
	profileHandler := profile.NewProfileHandler(profileService)
	invitationHandler := invitation.NewInvitationHandler(invitationService)
	customerHandler := customer.NewCustomerHandler(customerService)
	guestHandler := guest.NewGuestHandler(guestService)
	departmentHandler := department.NewDepartmentHandler(departmentService)
	auditHandler := audit.NewAuditHandler(auditService)
	unifiedHandler := unified.NewUnifiedHandler(unifiedService)
	integrationHandler := integration.NewIntegrationHandler(integrationService)
	adminHandler := admin.NewAdminHandler(adminService)

	// --- ROUTE REGISTRATION ---
	v1 := app.Group("/api/v1")

	// --- Public Routes ---
	authPublic := v1.Group("/auth")
	authPublic.Post("/register", authHandler.Register)
	authPublic.Post("/login", authHandler.Login)
	authPublic.Post("/refresh", authHandler.RefreshToken)
	authPublic.Post("/forgot-password", authHandler.ForgotPassword)
	authPublic.Post("/reset-password", authHandler.ResetPassword)

	v1.Post("/organizations", orgHandler.CreateOrganization)
	v1.Post("/customers/register", customerHandler.Register)
	v1.Post("/guests/session", guestHandler.CreateSession)
	v1.Post("/invitations/accept", invitationHandler.AcceptInvitation)

	// --- Protected Routes ---
	p := v1.Group("/", middleware.SmartAuthMiddleware())

	p.Post("/auth/logout", authHandler.Logout)
	p.Get("/auth/me", profileHandler.GetMyProfile)    // Added GET /me
	p.Put("/auth/me", profileHandler.UpdateMyProfile) // Added PUT /me

	// Profile routes
	profileRoutes := p.Group("/profiles")
	profileRoutes.Get("/me", profileHandler.GetMyProfile)
	profileRoutes.Put("/me", profileHandler.UpdateMyProfile)

	// Customer routes
	customerRoutes := p.Group("/customers")
	customerRoutes.Get("/profile", customerHandler.GetProfile)
	customerRoutes.Post("/addresses", customerHandler.AddAddress)
	customerRoutes.Get("/addresses", customerHandler.ListAddresses)
	customerRoutes.Put("/addresses/:id", customerHandler.UpdateAddress)
	customerRoutes.Delete("/addresses/:id", customerHandler.DeleteAddress)

	// Guest routes (some are public, some might be semi-protected by cookie)
	v1.Get("/guests/session", guestHandler.GetSession)
	v1.Put("/guests/session/cart", guestHandler.UpdateCart)

	// Organization-specific routes
	orgScoped := p.Group("/org/:slug")
	orgScoped.Delete("/", orgHandler.DeleteOrganization)
	orgScoped.Put("/", orgHandler.UpdateOrganization)

	userRoutes := orgScoped.Group("/users")
	userRoutes.Get("/", userHandler.ListUsers)
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Put("/:id/status", userHandler.UpdateUserStatus)

	roleRoutes := orgScoped.Group("/roles")
	roleRoutes.Post("/", roleHandler.CreateRole)
	roleRoutes.Get("/", roleHandler.ListRoles)
	roleRoutes.Put("/:id", roleHandler.UpdateRole)
	roleRoutes.Delete("/:id", roleHandler.DeleteRole)
	roleRoutes.Get("/:id/users", roleHandler.ListUsersInRole)

	deptRoutes := orgScoped.Group("/departments")
	deptRoutes.Post("/", departmentHandler.CreateDepartment)
	deptRoutes.Get("/", departmentHandler.ListDepartments)
	deptRoutes.Get("/:id", departmentHandler.GetDepartment)
	deptRoutes.Put("/:id", departmentHandler.UpdateDepartment)
	deptRoutes.Delete("/:id", departmentHandler.DeleteDepartment)

	invitationRoutes := orgScoped.Group("/invitations")
	invitationRoutes.Post("/", invitationHandler.CreateInvitation)
	invitationRoutes.Get("/", invitationHandler.ListInvitations)
	invitationRoutes.Delete("/:id", invitationHandler.CancelInvitation)

	// Unified, Integration, and Admin routes (placeholders)
	p.Get("/unified/dashboard", unifiedHandler.GetDashboard)
	p.Post("/integration/validate-token", integrationHandler.ValidateToken)
	p.Get("/admin/stats", adminHandler.GetSystemStats) // Super-admin only

	// Audit routes
	auditRoutes := orgScoped.Group("/audits")
	auditRoutes.Get("/", auditHandler.ListAuditLogs)
	// TODO: Add other audit routes

	// Add a placeholder for the top-level /users routes
	p.Get("/users", userHandler.ListUsers) // This would need to be adapted for non-org context
	p.Post("/users/bulk-create", userHandler.BulkCreateUsers)

	// Add missing customer routes
	customerRoutes.Post("/logout", authHandler.Logout)                              // Reuse auth logout
	customerRoutes.Post("/refresh", authHandler.RefreshToken)                       // Reuse auth refresh
	customerRoutes.Delete("/profile", customerHandler.DeleteProfile)                // Placeholder
	customerRoutes.Put("/addresses/:id/default", customerHandler.SetDefaultAddress) // Placeholder
	customerRoutes.Get("/preferences", customerHandler.GetPreferences)
	customerRoutes.Put("/preferences", customerHandler.UpdatePreferences)

	// Add missing org-scoped auth routes
	orgAuth := orgScoped.Group("/auth")
	orgAuth.Post("/login", authHandler.Login) // Can be handled by main login with org_slug
	// TODO: Add other org-auth routes

	// Health check routes
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
}
