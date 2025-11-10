package domain

import (
	"time"

	"github.com/google/uuid"
)

// AI Agent Types
type AgentType string

const (
	AgentTypeListing     AgentType = "listing"     // 出品エージェント
	AgentTypeNegotiation AgentType = "negotiation" // 交渉エージェント
	AgentTypeShipping    AgentType = "shipping"    // 配送エージェント
)

// Negotiation Mode - AI vs Manual
type NegotiationMode string

const (
	NegotiationModeAI     NegotiationMode = "ai"     // AI自動交渉
	NegotiationModeManual NegotiationMode = "manual" // 手動交渉
	NegotiationModeHybrid NegotiationMode = "hybrid" // AIアシスト付き手動
)

// AI Listing Agent - 出品時のAI生成データ
type AIListingData struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID            uuid.UUID `json:"product_id" gorm:"type:uuid;not null;uniqueIndex"`
	IsAIGenerated        bool      `json:"is_ai_generated" gorm:"default:true"`
	AIConfidenceScore    float64   `json:"ai_confidence_score" gorm:"type:decimal(5,2)"` // 0-100
	GeneratedTitle       string    `json:"generated_title"`
	GeneratedDescription string    `json:"generated_description"`
	GeneratedCategory    string    `json:"generated_category"`
	GeneratedCondition   string    `json:"generated_condition"`
	GeneratedPrice       int       `json:"generated_price"`
	UserModifiedFields   string    `json:"user_modified_fields" gorm:"type:text"`  // JSON array of field names
	ImageAnalysisResult  string    `json:"image_analysis_result" gorm:"type:text"` // AI画像分析結果
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// AI Negotiation Agent - オファー交渉のAI設定
type AINegotiationSettings struct {
	ID                   uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID            uuid.UUID       `json:"product_id" gorm:"type:uuid;not null;uniqueIndex"`
	Product              *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Mode                 NegotiationMode `json:"mode" gorm:"default:manual"` // ai, manual, hybrid
	IsEnabled            bool            `json:"is_enabled" gorm:"default:false"`
	MinAcceptablePrice   int             `json:"min_acceptable_price"`                  // AI交渉の最低価格
	AutoAcceptThreshold  int             `json:"auto_accept_threshold"`                 // この価格以上は自動承認
	AutoRejectThreshold  int             `json:"auto_reject_threshold"`                 // この価格未満は自動拒否
	NegotiationStrategy  string          `json:"negotiation_strategy"`                  // aggressive, moderate, conservative
	AIResponseTemplate   string          `json:"ai_response_template" gorm:"type:text"` // AIの返信テンプレート
	TotalOffersProcessed int             `json:"total_offers_processed" gorm:"default:0"`
	AIAcceptedCount      int             `json:"ai_accepted_count" gorm:"default:0"`
	AIRejectedCount      int             `json:"ai_rejected_count" gorm:"default:0"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// AI Shipping Preparation - 配送準備のAI提案
type AIShippingPreparation struct {
	ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PurchaseID           uuid.UUID  `json:"purchase_id" gorm:"type:uuid;not null;uniqueIndex"`
	Purchase             *Purchase  `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
	IsAIPrepared         bool       `json:"is_ai_prepared" gorm:"default:true"`
	SuggestedCarrier     string     `json:"suggested_carrier"`                         // 推奨配送業者
	SuggestedPackageSize string     `json:"suggested_package_size"`                    // 60サイズ、80サイズ等
	EstimatedWeight      float64    `json:"estimated_weight" gorm:"type:decimal(8,2)"` // kg
	EstimatedCost        int        `json:"estimated_cost"`                            // 配送料
	SuggestedLabel       string     `json:"suggested_label" gorm:"type:text"`          // 送り状情報
	TrackingNumber       string     `json:"tracking_number"`
	ShippingInstructions string     `json:"shipping_instructions" gorm:"type:text"` // 梱包指示
	UserApproved         bool       `json:"user_approved" gorm:"default:false"`
	UserModifications    string     `json:"user_modifications" gorm:"type:text"` // JSON
	ApprovedAt           *time.Time `json:"approved_at"`
	ShippedAt            *time.Time `json:"shipped_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// AI Agent Activity Log - エージェント活動ログ
type AIAgentLog struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	AgentType   AgentType `json:"agent_type" gorm:"not null;index"`
	Action      string    `json:"action"`                     // generated, accepted, rejected, modified
	TargetID    uuid.UUID `json:"target_id" gorm:"type:uuid"` // ProductID, OfferID, PurchaseID
	Details     string    `json:"details" gorm:"type:text"`   // JSON details
	Success     bool      `json:"success" gorm:"default:true"`
	ErrorMsg    string    `json:"error_msg"`
	ProcessTime int       `json:"process_time"` // milliseconds
	CreatedAt   time.Time `json:"created_at"`
}

// Repository Interface
type AIAgentRepository interface {
	// Listing Agent
	CreateListingData(data *AIListingData) error
	GetListingData(productID uuid.UUID) (*AIListingData, error)
	UpdateListingData(data *AIListingData) error

	// Negotiation Agent
	CreateNegotiationSettings(settings *AINegotiationSettings) error
	GetNegotiationSettings(productID uuid.UUID) (*AINegotiationSettings, error)
	UpdateNegotiationSettings(settings *AINegotiationSettings) error
	DeleteNegotiationSettings(productID uuid.UUID) error

	// Shipping Agent
	CreateShippingPreparation(prep *AIShippingPreparation) error
	GetShippingPreparation(purchaseID uuid.UUID) (*AIShippingPreparation, error)
	UpdateShippingPreparation(prep *AIShippingPreparation) error
	ApproveShipping(purchaseID uuid.UUID, modifications string) error

	// Logs
	CreateLog(log *AIAgentLog) error
	GetUserLogs(userID uuid.UUID, agentType AgentType, limit int) ([]*AIAgentLog, error)
	GetAgentStats(userID uuid.UUID) (*AIAgentStats, error)
}

// AI Agent Statistics
type AIAgentStats struct {
	TotalAIGenerations  int64   `json:"total_ai_generations"`
	ListingsCreated     int64   `json:"listings_created"`
	NegotiationsHandled int64   `json:"negotiations_handled"`
	ShipmentsPrepared   int64   `json:"shipments_prepared"`
	AverageConfidence   float64 `json:"average_confidence"`
	TimeSavedMinutes    int64   `json:"time_saved_minutes"`
	AcceptanceRate      float64 `json:"acceptance_rate"`
}

// Request/Response DTOs

// AI Listing Generation Request
type AIListingGenerationRequest struct {
	ImageURLs   []string `json:"image_urls" binding:"required"`
	UserHints   string   `json:"user_hints"`   // Optional user hints
	AutoPublish bool     `json:"auto_publish"` // 承認なしで公開するか
}

// AI Listing Generation Response
type AIListingGenerationResponse struct {
	ProductID           uuid.UUID             `json:"product_id"`
	ListingData         *AIListingData        `json:"listing_data"`
	SuggestedProduct    *SuggestedProductData `json:"suggested_product"`
	ConfidenceBreakdown map[string]float64    `json:"confidence_breakdown"`
	RequiresApproval    bool                  `json:"requires_approval"`
}

// Suggested Product Data for Approval Screen
type SuggestedProductData struct {
	Title             string           `json:"title"`
	Description       string           `json:"description"`
	Category          string           `json:"category"`
	Condition         ProductCondition `json:"condition"`
	Price             int              `json:"price"`
	WeightKg          float64          `json:"weight_kg"`
	DetectedBrand     string           `json:"detected_brand"`
	DetectedModel     string           `json:"detected_model"`
	KeyFeatures       []string         `json:"key_features"`
	PricingRationale  string           `json:"pricing_rationale"`
	CategoryRationale string           `json:"category_rationale"`
}

// AI Negotiation Toggle Request
type ToggleAINegotiationRequest struct {
	ProductID           uuid.UUID       `json:"product_id" binding:"required"`
	Mode                NegotiationMode `json:"mode" binding:"required"`
	MinAcceptablePrice  int             `json:"min_acceptable_price" binding:"required"`
	AutoAcceptThreshold int             `json:"auto_accept_threshold"`
	AutoRejectThreshold int             `json:"auto_reject_threshold"`
	Strategy            string          `json:"strategy"` // aggressive, moderate, conservative
}

// AI Shipping Preparation Request
type AIShippingPreparationRequest struct {
	PurchaseID uuid.UUID `json:"purchase_id" binding:"required"`
}

// AI Shipping Approval Request
type ApproveShippingRequest struct {
	Approved      bool   `json:"approved" binding:"required"`
	Carrier       string `json:"carrier"`
	PackageSize   string `json:"package_size"`
	Modifications string `json:"modifications"` // JSON string
}
