package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuctionStatus string

const (
	AuctionStatusActive    AuctionStatus = "active"
	AuctionStatusCompleted AuctionStatus = "completed"
	AuctionStatusCancelled AuctionStatus = "cancelled"
)

type Auction struct {
	ID              uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID       uuid.UUID     `json:"product_id" gorm:"type:uuid;not null;index"`
	Product         *Product      `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	SellerID        uuid.UUID     `json:"seller_id" gorm:"type:uuid;not null;index"`
	StartPrice      int           `json:"start_price" gorm:"not null"`
	CurrentBid      int           `json:"current_bid" gorm:"default:0"`
	MinBidIncrement int           `json:"min_bid_increment" gorm:"default:100"`
	WinnerID        *uuid.UUID    `json:"winner_id" gorm:"type:uuid"`
	Winner          *User         `json:"winner,omitempty" gorm:"foreignKey:WinnerID"`
	Status          AuctionStatus `json:"status" gorm:"default:active"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

type Bid struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AuctionID uuid.UUID `json:"auction_id" gorm:"type:uuid;not null;index"`
	Auction   *Auction  `json:"auction,omitempty" gorm:"foreignKey:AuctionID"`
	BidderID  uuid.UUID `json:"bidder_id" gorm:"type:uuid;not null;index"`
	Bidder    *User     `json:"bidder,omitempty" gorm:"foreignKey:BidderID"`
	Amount    int       `json:"amount" gorm:"not null"`
	IsWinning bool      `json:"is_winning" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

type AuctionRepository interface {
	Create(auction *Auction) error
	FindByID(id uuid.UUID) (*Auction, error)
	FindByProductID(productID uuid.UUID) (*Auction, error)
	FindActiveAuctions() ([]*Auction, error)
	Update(auction *Auction) error
}

type BidRepository interface {
	Create(bid *Bid) error
	FindByAuctionID(auctionID uuid.UUID) ([]*Bid, error)
	FindWinningBid(auctionID uuid.UUID) (*Bid, error)
	UpdateWinningStatus(auctionID, bidID uuid.UUID) error
}

type CreateAuctionRequest struct {
	ProductID       string `json:"product_id" binding:"required"`
	StartPrice      int    `json:"start_price" binding:"required,min=1"`
	MinBidIncrement int    `json:"min_bid_increment" binding:"required,min=1"`
	DurationMinutes int    `json:"duration_minutes" binding:"required,min=1,max=10080"` // Max 1 week
}

type PlaceBidRequest struct {
	Amount int `json:"amount" binding:"required,min=1"`
}
