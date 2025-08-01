package middleware

import (
	"errors"
	"practicev2/database"
	"practicev2/module/authentication/auth"
	"practicev2/module/authentication/models"
	"practicev2/module/authentication/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthContext holds the verified, rich context of the authenticated user for a given request.
type AuthContext struct {
	Identity        *models.Identity
	Memberships     []models.OrganizationalMembership
	CustomerProfile *models.CustomerProfile
	CurrentOrg      *models.Organization
	CurrentRole     *models.Role
	IsGuest         bool
}

// SmartAuthMiddleware is an intelligent middleware that validates a unified JWT
// and builds a rich authentication context.
func SmartAuthMiddleware() fiber.Handler {
	// Instantiate services needed by the middleware
	jwtService := utils.NewJWTService()
	// The repository is needed to fetch full context details
	repo := auth.NewAuthRepository(database.DBconn)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// This could be a guest user, allow proceeding but mark as guest.
			// For protected routes, a subsequent permission check will fail.
			c.Locals("authContext", &AuthContext{IsGuest: true})
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Malformed JWT")
		}
		tokenString := parts[1]

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid or expired JWT", err)
		}

		// Build the rich AuthContext from the claims
		orgSlug := c.Params("slug")

		// Si no se obtuvo de los parámetros, intentar extraer de la URL manualmente
		if orgSlug == "" {
			path := c.OriginalURL()
			// Formato esperado: /api/v1/org/{slug}/...
			if strings.Contains(path, "/org/") {
				parts := strings.Split(path, "/org/")
				if len(parts) > 1 {
					slugPart := strings.Split(parts[1], "/")[0]
					if slugPart != "" {
						orgSlug = slugPart
					}
				}
			}
		}

		authCtx, err := buildAuthContext(repo, claims, orgSlug)
		if err != nil {
			return utils.SendError(c, fiber.StatusForbidden, "Invalid authentication context", err)
		}

		c.Locals("authContext", authCtx)
		c.Locals("identity_id", claims.IdentityID)
		return c.Next()
	}
}

// buildAuthContext constructs the rich AuthContext from JWT claims and request data.
func buildAuthContext(repo auth.AuthRepository, claims *utils.UnifiedClaims, orgSlug string) (*AuthContext, error) {
	if claims.IdentityID == "" {
		return nil, errors.New("invalid token: missing identity_id")
	}

	identity, memberships, customer, err := repo.GetFullIdentityContext(claims.IdentityID)
	if err != nil {
		return nil, errors.New("failed to retrieve full identity context")
	}

	ctx := &AuthContext{
		Identity:        identity,
		Memberships:     memberships,
		CustomerProfile: customer,
		IsGuest:         false,
	}

	// If an organization context is specified in the request, find the relevant membership
	if orgSlug != "" {
		var currentMembership *models.OrganizationalMembership

		for _, m := range memberships {
			if m.Organization.Slug == orgSlug {
				currentMembership = &m
				break
			}
		}

		if currentMembership == nil {
			return nil, errors.New("user is not a member of the specified organization")
		}

		ctx.CurrentOrg = &currentMembership.Organization
		ctx.CurrentRole = &currentMembership.Role
	}

	return ctx, nil
}
