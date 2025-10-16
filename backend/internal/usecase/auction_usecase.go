package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type AuctionUseCase interface {
	CreateAuction(sellerID, productID uuid.UUID, startPrice, minIncrement, durationMinutes int) (*domain.Auction, error)
	PlaceBid(auctionID, bidderID uuid.UUID, amount int) (*domain.Bid, error)
	GetAuction(auctionID uuid.UUID) (*domain.Auction, error)
	GetActiveAuctions() ([]*domain.Auction, error)
	GetAuctionBids(auctionID uuid.UUID) ([]*domain.Bid, error)
	CompleteAuction(auctionID uuid.UUID) error
}

type auctionUseCase struct {
	auctionRepo domain.AuctionRepository
	bidRepo     domain.BidRepository
	productRepo domain.ProductRepository
}

func NewAuctionUseCase(
	auctionRepo domain.AuctionRepository,
	bidRepo domain.BidRepository,
	productRepo domain.ProductRepository,
) AuctionUseCase {
	return &auctionUseCase{
		auctionRepo: auctionRepo,
		bidRepo:     bidRepo,
		productRepo: productRepo,
	}
}

func (u *auctionUseCase) CreateAuction(sellerID, productID uuid.UUID, startPrice, minIncrement, durationMinutes int) (*domain.Auction, error) {
	// Get product
	product, err := u.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check ownership
	if product.SellerID != sellerID {
		return nil, errors.New("unauthorized: not product owner")
	}

	// Check if product is active
	if product.Status != domain.StatusActive {
		return nil, errors.New("product is not available")
	}

	// Create auction
	now := time.Now()
	auction := &domain.Auction{
		ProductID:       productID,
		SellerID:        sellerID,
		StartPrice:      startPrice,
		CurrentBid:      startPrice,
		MinBidIncrement: minIncrement,
		Status:          domain.AuctionStatusActive,
		StartTime:       now,
		EndTime:         now.Add(time.Duration(durationMinutes) * time.Minute),
	}

	if err := u.auctionRepo.Create(auction); err != nil {
		return nil, err
	}

	return u.auctionRepo.FindByID(auction.ID)
}

func (u *auctionUseCase) PlaceBid(auctionID, bidderID uuid.UUID, amount int) (*domain.Bid, error) {
	// Get auction
	auction, err := u.auctionRepo.FindByID(auctionID)
	if err != nil {
		return nil, errors.New("auction not found")
	}

	// Check if auction is active
	if auction.Status != domain.AuctionStatusActive {
		return nil, errors.New("auction is not active")
	}

	// Check if auction has ended
	if time.Now().After(auction.EndTime) {
		return nil, errors.New("auction has ended")
	}

	// Cannot bid on own auction
	if auction.SellerID == bidderID {
		return nil, errors.New("cannot bid on your own auction")
	}

	// Check minimum bid
	minBid := auction.CurrentBid + auction.MinBidIncrement
	if amount < minBid {
		return nil, errors.New("bid amount must be at least " + string(rune(minBid)))
	}

	// Create bid
	bid := &domain.Bid{
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    amount,
		IsWinning: true,
	}

	if err := u.bidRepo.Create(bid); err != nil {
		return nil, err
	}

	// Update auction current bid and winning status
	if err := u.bidRepo.UpdateWinningStatus(auctionID, bid.ID); err != nil {
		return nil, err
	}

	auction.CurrentBid = amount
	if err := u.auctionRepo.Update(auction); err != nil {
		return nil, err
	}

	return bid, nil
}

func (u *auctionUseCase) GetAuction(auctionID uuid.UUID) (*domain.Auction, error) {
	return u.auctionRepo.FindByID(auctionID)
}

func (u *auctionUseCase) GetActiveAuctions() ([]*domain.Auction, error) {
	return u.auctionRepo.FindActiveAuctions()
}

func (u *auctionUseCase) GetAuctionBids(auctionID uuid.UUID) ([]*domain.Bid, error) {
	return u.bidRepo.FindByAuctionID(auctionID)
}

func (u *auctionUseCase) CompleteAuction(auctionID uuid.UUID) error {
	auction, err := u.auctionRepo.FindByID(auctionID)
	if err != nil {
		return errors.New("auction not found")
	}

	// Check if auction has ended
	if time.Now().Before(auction.EndTime) {
		return errors.New("auction has not ended yet")
	}

	// Get winning bid
	winningBid, err := u.bidRepo.FindWinningBid(auctionID)
	if err != nil {
		// No bids, just mark as completed
		auction.Status = domain.AuctionStatusCompleted
		return u.auctionRepo.Update(auction)
	}

	// Update auction with winner
	auction.WinnerID = &winningBid.BidderID
	auction.Status = domain.AuctionStatusCompleted

	return u.auctionRepo.Update(auction)
}
