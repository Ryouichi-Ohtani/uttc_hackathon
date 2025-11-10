package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	DisplayName  string    `json:"display_name"`
	AvatarURL    string    `json:"avatar_url"`
	Bio          string    `json:"bio"`
	Role         UserRole  `json:"role" gorm:"default:'user'"`
	// Address information for shipping
	PostalCode   string     `json:"postal_code"`
	Prefecture   string     `json:"prefecture"`
	City         string     `json:"city"`
	AddressLine1 string     `json:"address_line1"`
	AddressLine2 string     `json:"address_line2"`
	PhoneNumber  string     `json:"phone_number"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-" gorm:"index"`
}

func (u *User) HasRole(role UserRole) bool {
	switch role {
	case RoleAdmin:
		return u.Role == RoleAdmin
	case RoleModerator:
		return u.Role == RoleAdmin || u.Role == RoleModerator
	case RoleUser:
		return true
	}
	return false
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
	UpdateSustainabilityStats(userID uuid.UUID, co2SavedKg float64) error
	GetLeaderboard(limit int, period string) ([]*LeaderboardEntry, error)
}

type LeaderboardEntry struct {
	Rank                int     `json:"rank"`
	User                *User   `json:"user"`
	TotalCO2SavedKg     float64 `json:"total_co2_saved_kg"`
	SustainabilityScore int     `json:"sustainability_score"`
	Level               int     `json:"level"`
}

// RegisterRequest represents user registration input
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=30"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"omitempty,min=1,max=100"`
}

// LoginRequest represents user login input
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
