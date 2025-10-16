package infrastructure

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
	"gorm.io/gorm"
)

type auctionRepository struct {
	db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) domain.AuctionRepository {
	return &auctionRepository{db: db}
}

func (r *auctionRepository) Create(auction *domain.Auction) error {
	return r.db.Create(auction).Error
}

func (r *auctionRepository) FindByID(id uuid.UUID) (*domain.Auction, error) {
	var auction domain.Auction
	err := r.db.Preload("Product").Preload("Product.Images").Preload("Winner").
		First(&auction, "id = ?", id).Error
	return &auction, err
}

func (r *auctionRepository) FindByProductID(productID uuid.UUID) (*domain.Auction, error) {
	var auction domain.Auction
	err := r.db.Preload("Product").Preload("Winner").
		Where("product_id = ? AND status = ?", productID, domain.AuctionStatusActive).
		First(&auction).Error
	return &auction, err
}

func (r *auctionRepository) FindActiveAuctions() ([]*domain.Auction, error) {
	var auctions []*domain.Auction
	err := r.db.Preload("Product").Preload("Product.Images").
		Where("status = ? AND end_time > ?", domain.AuctionStatusActive, time.Now()).
		Order("end_time ASC").
		Find(&auctions).Error
	return auctions, err
}

func (r *auctionRepository) Update(auction *domain.Auction) error {
	return r.db.Save(auction).Error
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) domain.BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) Create(bid *domain.Bid) error {
	return r.db.Create(bid).Error
}

func (r *bidRepository) FindByAuctionID(auctionID uuid.UUID) ([]*domain.Bid, error) {
	var bids []*domain.Bid
	err := r.db.Preload("Bidder").
		Where("auction_id = ?", auctionID).
		Order("created_at DESC").
		Find(&bids).Error
	return bids, err
}

func (r *bidRepository) FindWinningBid(auctionID uuid.UUID) (*domain.Bid, error) {
	var bid domain.Bid
	err := r.db.Preload("Bidder").
		Where("auction_id = ? AND is_winning = ?", auctionID, true).
		First(&bid).Error
	return &bid, err
}

func (r *bidRepository) UpdateWinningStatus(auctionID, bidID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Reset all bids for this auction
		if err := tx.Model(&domain.Bid{}).
			Where("auction_id = ?", auctionID).
			Update("is_winning", false).Error; err != nil {
			return err
		}

		// Set the new winning bid
		return tx.Model(&domain.Bid{}).
			Where("id = ?", bidID).
			Update("is_winning", true).Error
	})
}
