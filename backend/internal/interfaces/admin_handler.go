package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AdminHandler struct {
	productUseCase *usecase.ProductUseCase
	authUseCase    *usecase.AuthUseCase
}

func NewAdminHandler(productUseCase *usecase.ProductUseCase, authUseCase *usecase.AuthUseCase) *AdminHandler {
	return &AdminHandler{
		productUseCase: productUseCase,
		authUseCase:    authUseCase,
	}
}

// GetAllProducts returns all products for admin (no seller filter)
func (h *AdminHandler) GetAllProducts(c *gin.Context) {
	filters := domain.ProductFilters{
		Page:  1,
		Limit: 100,
	}

	products, pagination, err := h.productUseCase.List(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch products",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"pagination": pagination,
	})
}

// AdminUpdateProduct allows admin to update any product
func (h *AdminHandler) AdminUpdateProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid product ID",
			},
		})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	// Get existing product
	product, err := h.productUseCase.GetByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Product not found",
			},
		})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Condition != nil {
		updates["condition"] = *req.Condition
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	product, err = h.productUseCase.Update(productID, product.SellerID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update product",
			},
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// AdminDeleteProduct allows admin to delete any product
func (h *AdminHandler) AdminDeleteProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid product ID",
			},
		})
		return
	}

	// Get product to find seller ID
	product, err := h.productUseCase.GetByID(productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Product not found",
			},
		})
		return
	}

	if err := h.productUseCase.Delete(productID, product.SellerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DELETE_FAILED",
				"message": "Failed to delete product",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetAllUsers returns all users for admin
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	// This would need a UserRepository.List method
	// For now, return a simple message
	c.JSON(http.StatusOK, gin.H{"message": "User list endpoint"})
}

type UpdateProductRequest struct {
	Title       *string                  `json:"title"`
	Description *string                  `json:"description"`
	Price       *int                     `json:"price"`
	Category    *string                  `json:"category"`
	Condition   *domain.ProductCondition `json:"condition"`
	Status      *domain.ProductStatus    `json:"status"`
}
