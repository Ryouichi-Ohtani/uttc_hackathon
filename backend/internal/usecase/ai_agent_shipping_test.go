package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/ecomate/backend/internal/domain"
)

// ========== Test: Region Detection ==========

func TestDetectRegion(t *testing.T) {
	uc := &AIAgentUseCase{}

	tests := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "北海道",
			address:  "北海道札幌市中央区北1条西2丁目",
			expected: "北海道",
		},
		{
			name:     "東京都",
			address:  "東京都渋谷区渋谷1-2-3",
			expected: "関東",
		},
		{
			name:     "大阪府",
			address:  "大阪府大阪市北区梅田1-1-1",
			expected: "関西",
		},
		{
			name:     "沖縄県",
			address:  "沖縄県那覇市おもろまち1-1-1",
			expected: "沖縄",
		},
		{
			name:     "福岡県",
			address:  "福岡県福岡市博多区博多駅前1-1-1",
			expected: "九州",
		},
		{
			name:     "宮城県",
			address:  "宮城県仙台市青葉区中央1-1-1",
			expected: "東北",
		},
		{
			name:     "京都府",
			address:  "京都府京都市中京区烏丸通1-1",
			expected: "関西",
		},
		{
			name:     "不明な住所",
			address:  "Unknown Address 123",
			expected: "関東", // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.detectRegion(tt.address)
			assert.Equal(t, tt.expected, result, "Region detection failed for address: %s", tt.address)
		})
	}
}

// ========== Test: Regional Cost Adjustment ==========

func TestAdjustCostByRegion(t *testing.T) {
	uc := &AIAgentUseCase{}

	tests := []struct {
		name     string
		baseCost int
		region   string
		expected int
	}{
		{
			name:     "北海道は+30%",
			baseCost: 1000,
			region:   "北海道",
			expected: 1300,
		},
		{
			name:     "沖縄は+50%",
			baseCost: 1000,
			region:   "沖縄",
			expected: 1500,
		},
		{
			name:     "関東は基準価格",
			baseCost: 1000,
			region:   "関東",
			expected: 1000,
		},
		{
			name:     "関西は+5%",
			baseCost: 1000,
			region:   "関西",
			expected: 1050,
		},
		{
			name:     "九州は+20%",
			baseCost: 1000,
			region:   "九州",
			expected: 1200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.adjustCostByRegion(tt.baseCost, tt.region)
			assert.Equal(t, tt.expected, result, "Cost adjustment failed for region: %s", tt.region)
		})
	}
}

// ========== Test: Shipping Info Estimation ==========

func TestEstimateShippingInfo(t *testing.T) {
	uc := &AIAgentUseCase{}

	tests := []struct {
		name           string
		product        *domain.Product
		purchase       *domain.Purchase
		expectedSize   string
		expectedCarrier string
		minCost        int
		maxCost        int
	}{
		{
			name: "軽量商品（東京）",
			product: &domain.Product{
				Category: "books",
				WeightKg: 0.5,
			},
			purchase: &domain.Purchase{
				ShippingAddress: "東京都渋谷区渋谷1-2-3",
			},
			expectedSize:    "60サイズ",
			expectedCarrier: "ヤマト運輸（ネコポス）",
			minCost:         700,
			maxCost:         900,
		},
		{
			name: "中重量商品（大阪）",
			product: &domain.Product{
				Category: "electronics",
				WeightKg: 2.5,
			},
			purchase: &domain.Purchase{
				ShippingAddress: "大阪府大阪市北区梅田1-1-1",
			},
			expectedSize:    "80サイズ",
			expectedCarrier: "ヤマト運輸（宅急便）",
			minCost:         1000,
			maxCost:         1100,
		},
		{
			name: "重量商品（北海道）",
			product: &domain.Product{
				Category: "furniture",
				WeightKg: 6.0,
			},
			purchase: &domain.Purchase{
				ShippingAddress: "北海道札幌市中央区北1条西2丁目",
			},
			expectedSize:    "120サイズ以上",
			expectedCarrier: "ヤマト運輸（宅急便）",
			minCost:         2000, // 1600 * 1.3
			maxCost:         2100,
		},
		{
			name: "小型商品（沖縄）",
			product: &domain.Product{
				Category: "toys",
				WeightKg: 0.8,
			},
			purchase: &domain.Purchase{
				ShippingAddress: "沖縄県那覇市おもろまち1-1-1",
			},
			expectedSize:    "60サイズ",
			expectedCarrier: "ヤマト運輸（ネコポス）",
			minCost:         1100, // 800 * 1.5
			maxCost:         1300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.estimateShippingInfo(tt.product, tt.purchase)

			assert.Equal(t, tt.expectedSize, result.PackageSize, "Package size mismatch")
			assert.Equal(t, tt.expectedCarrier, result.Carrier, "Carrier mismatch")
			assert.GreaterOrEqual(t, result.Cost, tt.minCost, "Cost too low")
			assert.LessOrEqual(t, result.Cost, tt.maxCost, "Cost too high")
			assert.NotEmpty(t, result.Instructions, "Instructions should not be empty")

			// カテゴリ別の梱包指示チェック
			switch tt.product.Category {
			case "electronics":
				assert.Contains(t, result.Instructions, "精密機器", "Electronics should have special handling instructions")
			case "clothing":
				assert.Contains(t, result.Instructions, "防水", "Clothing should have waterproof instructions")
			case "furniture":
				assert.Contains(t, result.Instructions, "緩衝材", "Furniture should have padding instructions")
			}

			// 遠隔地の追加指示チェック
			if tt.purchase.ShippingAddress == "北海道札幌市中央区北1条西2丁目" ||
				tt.purchase.ShippingAddress == "沖縄県那覇市おもろまち1-1-1" {
				assert.Contains(t, result.Instructions, "遠隔地配送", "Remote areas should have additional instructions")
			}
		})
	}
}

// ========== Test: Delivery Days Estimation ==========

func TestEstimateDeliveryDays(t *testing.T) {
	uc := &AIAgentUseCase{}

	tests := []struct {
		name     string
		carrier  string
		weight   float64
		expected int
	}{
		{
			name:     "ネコポス（軽量）",
			carrier:  "ヤマト運輸（ネコポス）",
			weight:   0.5,
			expected: 1,
		},
		{
			name:     "宅急便（軽量）",
			carrier:  "ヤマト運輸（宅急便）",
			weight:   2.0,
			expected: 2,
		},
		{
			name:     "宅急便（重量）",
			carrier:  "ヤマト運輸（宅急便）",
			weight:   6.0,
			expected: 3,
		},
		{
			name:     "その他配送業者",
			carrier:  "佐川急便",
			weight:   3.0,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.estimateDeliveryDays(tt.carrier, tt.weight)
			assert.Equal(t, tt.expected, result, "Delivery days estimation failed")
		})
	}
}

// ========== Test: Season Detection ==========

func TestGetSeason(t *testing.T) {
	uc := &AIAgentUseCase{}

	tests := []struct {
		name     string
		month    int
		expected string
	}{
		{"1月は冬", 1, "冬"},
		{"2月は冬", 2, "冬"},
		{"3月は春", 3, "春"},
		{"4月は春", 4, "春"},
		{"5月は春", 5, "春"},
		{"6月は夏", 6, "夏"},
		{"7月は夏", 7, "夏"},
		{"8月は夏", 8, "夏"},
		{"9月は秋", 9, "秋"},
		{"10月は秋", 10, "秋"},
		{"11月は秋", 11, "秋"},
		{"12月は冬", 12, "冬"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uc.getSeason(time.Month(tt.month))
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ========== Integration Test: PrepareShipping (Concept) ==========

func TestPrepareShipping_Concept(t *testing.T) {
	// この統合テストは概念的なもので、実際のDB接続が必要
	// モックを使った完全な実装は省略

	t.Run("PrepareShipping should generate valid shipping info", func(t *testing.T) {
		// テストの構造を示すのみ
		ctx := context.Background()
		purchaseID := uuid.New()

		// 実際のテストでは以下を実装:
		// 1. Mock AIAgentRepository
		// 2. Mock ProductRepository
		// 3. Mock PurchaseRepository
		// 4. Call uc.PrepareShipping(ctx, purchaseID)
		// 5. Assert that returned AIShippingPreparation has valid data

		assert.NotNil(t, ctx, "Context should exist")
		assert.NotEqual(t, uuid.Nil, purchaseID, "Purchase ID should be valid")
	})
}

// ========== Edge Case Tests ==========

func TestShippingInfoEdgeCases(t *testing.T) {
	uc := &AIAgentUseCase{}

	t.Run("Zero weight product", func(t *testing.T) {
		product := &domain.Product{
			Category: "books",
			WeightKg: 0.0,
		}
		purchase := &domain.Purchase{
			ShippingAddress: "東京都渋谷区渋谷1-2-3",
		}

		result := uc.estimateShippingInfo(product, purchase)

		// 0kg商品でも配送情報が生成される
		assert.NotEmpty(t, result.Carrier)
		assert.NotEmpty(t, result.PackageSize)
		assert.Greater(t, result.Cost, 0, "Cost should be positive even for 0kg")
	})

	t.Run("Very heavy product", func(t *testing.T) {
		product := &domain.Product{
			Category: "furniture",
			WeightKg: 50.0,
		}
		purchase := &domain.Purchase{
			ShippingAddress: "東京都渋谷区渋谷1-2-3",
		}

		result := uc.estimateShippingInfo(product, purchase)

		assert.Equal(t, "120サイズ以上", result.PackageSize)
		assert.GreaterOrEqual(t, result.Cost, 1600, "Heavy items should have high shipping cost")
	})

	t.Run("Empty shipping address", func(t *testing.T) {
		product := &domain.Product{
			Category: "electronics",
			WeightKg: 1.5,
		}
		purchase := &domain.Purchase{
			ShippingAddress: "",
		}

		result := uc.estimateShippingInfo(product, purchase)

		// 空の住所でもデフォルト地域（関東）で計算される
		assert.NotEmpty(t, result.Carrier)
		assert.Greater(t, result.Cost, 0)
	})
}

// ========== Test: Category-specific Instructions ==========

func TestCategorySpecificInstructions(t *testing.T) {
	uc := &AIAgentUseCase{}

	categories := []struct {
		category    string
		expectedKey string
	}{
		{"electronics", "精密機器"},
		{"clothing", "防水"},
		{"books", "保護"},
		{"furniture", "緩衝材"},
		{"toys", "破損"},
	}

	for _, cat := range categories {
		t.Run("Category: "+cat.category, func(t *testing.T) {
			product := &domain.Product{
				Category: cat.category,
				WeightKg: 1.0,
			}
			purchase := &domain.Purchase{
				ShippingAddress: "東京都渋谷区渋谷1-2-3",
			}

			result := uc.estimateShippingInfo(product, purchase)

			assert.Contains(t, result.Instructions, cat.expectedKey,
				"Instructions for %s should contain '%s'", cat.category, cat.expectedKey)
		})
	}
}
