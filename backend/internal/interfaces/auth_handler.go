package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("VALIDATION_ERROR", err.Error(), nil))
		return
	}

	resp, err := h.authUseCase.Register(&req)
	if err != nil {
		c.JSON(http.StatusConflict, ErrorResponse("REGISTRATION_FAILED", err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("VALIDATION_ERROR", err.Error(), nil))
		return
	}

	resp, err := h.authUseCase.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse("LOGIN_FAILED", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := GetUserIDFromContext(c)

	user, err := h.authUseCase.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse("USER_NOT_FOUND", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, user)
}
