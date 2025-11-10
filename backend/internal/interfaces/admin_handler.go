package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

type AdminHandler struct {
	productUseCase  *usecase.ProductUseCase
	authUseCase     *usecase.AuthUseCase
	purchaseUseCase usecase.PurchaseUseCase
}

func NewAdminHandler(productUseCase *usecase.ProductUseCase, authUseCase *usecase.AuthUseCase, purchaseUseCase usecase.PurchaseUseCase) *AdminHandler {
	return &AdminHandler{
		productUseCase:  productUseCase,
		authUseCase:     authUseCase,
		purchaseUseCase: purchaseUseCase,
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
	page := 1
	limit := 100

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := parseInt(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := parseInt(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	users, total, err := h.authUseCase.ListUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch users",
			},
		})
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// UpdateUser allows admin to update any user
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid user ID",
			},
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.PostalCode != nil {
		updates["postal_code"] = *req.PostalCode
	}
	if req.Prefecture != nil {
		updates["prefecture"] = *req.Prefecture
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.AddressLine1 != nil {
		updates["address_line1"] = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		updates["address_line2"] = *req.AddressLine2
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}

	user, err := h.authUseCase.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser allows admin to delete any user
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid user ID",
			},
		})
		return
	}

	if err := h.authUseCase.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "DELETE_FAILED",
				"message": "Failed to delete user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetAllPurchases returns all purchases for admin
func (h *AdminHandler) GetAllPurchases(c *gin.Context) {
	page := 1
	limit := 100

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := parseInt(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := parseInt(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	purchases, pagination, err := h.purchaseUseCase.List(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch purchases",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases":  purchases,
		"pagination": pagination,
	})
}

// UpdatePurchase allows admin to update purchase status
func (h *AdminHandler) UpdatePurchase(c *gin.Context) {
	purchaseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "Invalid purchase ID",
			},
		})
		return
	}

	var req UpdatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.purchaseUseCase.UpdatePurchaseStatus(purchaseID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update purchase",
			},
		})
		return
	}

	purchase, err := h.purchaseUseCase.GetByID(purchaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch updated purchase",
			},
		})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

// Helper function
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

type UpdateProductRequest struct {
	Title       *string                  `json:"title"`
	Description *string                  `json:"description"`
	Price       *int                     `json:"price"`
	Category    *string                  `json:"category"`
	Condition   *domain.ProductCondition `json:"condition"`
	Status      *domain.ProductStatus    `json:"status"`
}

type UpdateUserRequest struct {
	DisplayName  *string          `json:"display_name"`
	Role         *domain.UserRole `json:"role"`
	Bio          *string          `json:"bio"`
	PostalCode   *string          `json:"postal_code"`
	Prefecture   *string          `json:"prefecture"`
	City         *string          `json:"city"`
	AddressLine1 *string          `json:"address_line1"`
	AddressLine2 *string          `json:"address_line2"`
	PhoneNumber  *string          `json:"phone_number"`
}

type UpdatePurchaseRequest struct {
	Status domain.PurchaseStatus `json:"status" binding:"required"`
}
