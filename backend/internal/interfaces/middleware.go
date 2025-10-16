package interfaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

const UserIDKey = "user_id"

func AuthMiddleware(authUseCase *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse("UNAUTHORIZED", "Missing authorization header", nil))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrorResponse("UNAUTHORIZED", "Invalid authorization header format", nil))
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := authUseCase.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse("UNAUTHORIZED", "Invalid or expired token", nil))
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.Nil
	}
	return userID.(uuid.UUID)
}

func ErrorResponse(code, message string, details interface{}) gin.H {
	return gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"details": details,
		},
	}
}
