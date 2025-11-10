package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtSecret string, jwtExpiry int) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (uc *AuthUseCase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if email already exists
	if _, err := uc.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	if _, err := uc.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		DisplayName:  req.DisplayName,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (uc *AuthUseCase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (uc *AuthUseCase) GetUserByID(id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.FindByID(id)
}

func (uc *AuthUseCase) ListUsers(page, limit int) ([]*domain.User, int64, error) {
	return uc.userRepo.List(page, limit)
}

func (uc *AuthUseCase) UpdateUser(id uuid.UUID, updates map[string]interface{}) (*domain.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if displayName, ok := updates["display_name"].(string); ok {
		user.DisplayName = displayName
	}
	if role, ok := updates["role"].(domain.UserRole); ok {
		user.Role = role
	}
	if bio, ok := updates["bio"].(string); ok {
		user.Bio = bio
	}
	if postalCode, ok := updates["postal_code"].(string); ok {
		user.PostalCode = postalCode
	}
	if prefecture, ok := updates["prefecture"].(string); ok {
		user.Prefecture = prefecture
	}
	if city, ok := updates["city"].(string); ok {
		user.City = city
	}
	if addressLine1, ok := updates["address_line1"].(string); ok {
		user.AddressLine1 = addressLine1
	}
	if addressLine2, ok := updates["address_line2"].(string); ok {
		user.AddressLine2 = addressLine2
	}
	if phoneNumber, ok := updates["phone_number"].(string); ok {
		user.PhoneNumber = phoneNumber
	}

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) DeleteUser(id uuid.UUID) error {
	return uc.userRepo.Delete(id)
}

func (uc *AuthUseCase) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * time.Duration(uc.jwtExpiry)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid user_id in token")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user_id format: %w", err)
		}

		return userID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}
