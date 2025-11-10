package usecase

import (
	"strings"

	"github.com/yourusername/ecomate/backend/internal/domain"
)

type VoiceSearchUseCase interface {
	SearchByVoice(transcript string) ([]*domain.Product, error)
	ProcessVoiceQuery(audioData []byte) (*VoiceSearchResult, error)
}

type VoiceSearchResult struct {
	Transcript string            `json:"transcript"`
	Products   []*domain.Product `json:"products"`
	Intent     string            `json:"intent"`   // "search", "filter", "navigate"
	Entities   map[string]string `json:"entities"` // category, price_range, condition
}

type voiceSearchUseCase struct {
	productRepo domain.ProductRepository
}

func NewVoiceSearchUseCase(productRepo domain.ProductRepository) VoiceSearchUseCase {
	return &voiceSearchUseCase{
		productRepo: productRepo,
	}
}

func (u *voiceSearchUseCase) SearchByVoice(transcript string) ([]*domain.Product, error) {
	// Parse voice command
	// In production, use NLP/NLU service
	filters := u.parseVoiceCommand(transcript)

	products, _, err := u.productRepo.List(filters)
	return products, err
}

func (u *voiceSearchUseCase) ProcessVoiceQuery(audioData []byte) (*VoiceSearchResult, error) {
	// In production, use Google Speech-to-Text API or similar
	// For now, return mock result
	transcript := "search for electronics under 10000 yen"

	products, _ := u.SearchByVoice(transcript)

	entities := u.extractEntities(transcript)

	return &VoiceSearchResult{
		Transcript: transcript,
		Products:   products,
		Intent:     "search",
		Entities:   entities,
	}, nil
}

func (u *voiceSearchUseCase) parseVoiceCommand(transcript string) *domain.ProductFilters {
	filters := &domain.ProductFilters{
		Limit: 20,
	}

	lowerTranscript := strings.ToLower(transcript)

	// Extract category
	categories := []string{"electronics", "fashion", "home", "sports", "books"}
	for _, cat := range categories {
		if strings.Contains(lowerTranscript, cat) {
			filters.Category = cat
			break
		}
	}

	// Extract condition
	if strings.Contains(lowerTranscript, "new") {
		filters.Condition = domain.ConditionNew
	} else if strings.Contains(lowerTranscript, "like new") {
		filters.Condition = domain.ConditionLikeNew
	} else if strings.Contains(lowerTranscript, "good") {
		filters.Condition = domain.ConditionGood
	}

	// Extract price range
	if strings.Contains(lowerTranscript, "under") || strings.Contains(lowerTranscript, "less than") {
		// Parse price: "under 10000 yen"
		filters.MaxPrice = 10000 // Simplified
	}

	// Extract search keywords
	keywords := []string{"search for", "find", "looking for", "show me"}
	for _, keyword := range keywords {
		if idx := strings.Index(lowerTranscript, keyword); idx != -1 {
			searchText := lowerTranscript[idx+len(keyword):]
			// Clean up and extract main search term
			filters.Search = strings.TrimSpace(strings.Split(searchText, " ")[0])
			break
		}
	}

	return filters
}

func (u *voiceSearchUseCase) extractEntities(transcript string) map[string]string {
	entities := make(map[string]string)

	lowerTranscript := strings.ToLower(transcript)

	// Extract category
	categories := []string{"electronics", "fashion", "home", "sports", "books"}
	for _, cat := range categories {
		if strings.Contains(lowerTranscript, cat) {
			entities["category"] = cat
			break
		}
	}

	// Extract price range
	if strings.Contains(lowerTranscript, "under") {
		entities["price_range"] = "under_specified"
	}

	// Extract condition
	conditions := []string{"new", "like new", "good", "fair"}
	for _, cond := range conditions {
		if strings.Contains(lowerTranscript, cond) {
			entities["condition"] = cond
			break
		}
	}

	return entities
}
