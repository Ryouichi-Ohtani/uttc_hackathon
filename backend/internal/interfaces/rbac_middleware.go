package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

// RequireRole creates middleware that checks if user has required role
func RequireRole(authUseCase *usecase.AuthUseCase, requiredRole domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := authUseCase.GetUserByID(userID.(uuid.UUID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.HasRole(requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("user_role", user.Role)
		c.Next()
	}
}

// RequireAdmin is a shorthand for RequireRole(RoleAdmin)
func RequireAdmin(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return RequireRole(authUseCase, domain.RoleAdmin)
}

// RequireModerator is a shorthand for RequireRole(RoleModerator)
func RequireModerator(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return RequireRole(authUseCase, domain.RoleModerator)
}
