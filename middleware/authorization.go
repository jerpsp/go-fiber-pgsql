package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// UserRole type for roles
type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
)

// RoleAuthorization middleware checks if the user has the required role
func RoleAuthorization(allowedRoles ...UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user role from the context (set by JWTMiddleware)
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: role information not found",
			})
		}

		// Check if user has one of the allowed roles
		hasAllowedRole := false
		userRole := UserRole(role)

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				hasAllowedRole = true
				break
			}
		}

		if !hasAllowedRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions",
			})
		}

		return c.Next()
	}
}

// AdminOnly middleware restricts access to admin users only
func AdminOnly() fiber.Handler {
	return RoleAuthorization(RoleAdmin)
}

// ModeratorOrAdmin middleware restricts access to moderators and admins
func ModeratorOrAdmin() fiber.Handler {
	return RoleAuthorization(RoleAdmin, RoleModerator)
}
