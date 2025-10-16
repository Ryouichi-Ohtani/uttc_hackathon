package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

type SocialUseCase struct {
	followRepo       domain.FollowRepository
	shareRepo        domain.ProductShareRepository
	feedRepo         domain.UserFeedRepository
}

func NewSocialUseCase(
	followRepo domain.FollowRepository,
	shareRepo domain.ProductShareRepository,
	feedRepo domain.UserFeedRepository,
) *SocialUseCase {
	return &SocialUseCase{
		followRepo: followRepo,
		shareRepo:  shareRepo,
		feedRepo:   feedRepo,
	}
}

func (uc *SocialUseCase) FollowUser(followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}

	isFollowing, err := uc.followRepo.IsFollowing(followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return fmt.Errorf("already following")
	}

	follow := &domain.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := uc.followRepo.Create(follow); err != nil {
		return err
	}

	// Create feed for the followed user
	_ = uc.feedRepo.Create(&domain.UserFeed{
		UserID:      followingID,
		ActorID:     followerID,
		ActionType:  "followed",
		TargetID:    followingID,
		TargetType:  "user",
		Description: "さんがあなたをフォローしました",
	})

	return nil
}

func (uc *SocialUseCase) UnfollowUser(followerID, followingID uuid.UUID) error {
	return uc.followRepo.Delete(followerID, followingID)
}

func (uc *SocialUseCase) GetFollowers(userID uuid.UUID, limit int) ([]*domain.Follow, error) {
	return uc.followRepo.GetFollowers(userID, limit)
}

func (uc *SocialUseCase) GetFollowing(userID uuid.UUID, limit int) ([]*domain.Follow, error) {
	return uc.followRepo.GetFollowing(userID, limit)
}

func (uc *SocialUseCase) GetFollowStats(userID uuid.UUID) (followers, following int64, err error) {
	followers, err = uc.followRepo.GetFollowerCount(userID)
	if err != nil {
		return 0, 0, err
	}

	following, err = uc.followRepo.GetFollowingCount(userID)
	if err != nil {
		return 0, 0, err
	}

	return followers, following, nil
}

func (uc *SocialUseCase) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	return uc.followRepo.IsFollowing(followerID, followingID)
}

func (uc *SocialUseCase) ShareProduct(userID, productID uuid.UUID, platform string) (*domain.ProductShare, error) {
	shareURL := fmt.Sprintf("https://ecomate.app/products/%s?ref=%s", productID, userID)

	share := &domain.ProductShare{
		UserID:    userID,
		ProductID: productID,
		Platform:  platform,
		ShareURL:  shareURL,
	}

	if err := uc.shareRepo.Create(share); err != nil {
		return nil, err
	}

	// Create feed for followers
	_ = uc.feedRepo.CreateFeedForFollowers(
		userID,
		"shared",
		productID,
		"product",
		"さんが商品をシェアしました",
	)

	return share, nil
}

func (uc *SocialUseCase) GetFeed(userID uuid.UUID, limit int) ([]*domain.UserFeed, error) {
	return uc.feedRepo.GetFeed(userID, limit)
}
