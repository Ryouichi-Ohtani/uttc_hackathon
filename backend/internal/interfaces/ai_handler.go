package interfaces

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	geminiAPIKey string
}

func NewAIHandler() *AIHandler {
	return &AIHandler{
		geminiAPIKey: os.Getenv("GOOGLE_API_KEY"),
	}
}

type TranslateSearchRequest struct {
	Query string `json:"query" binding:"required"`
}

type TranslateSearchResponse struct {
	Japanese         string   `json:"japanese"`
	English          string   `json:"english"`
	Romanized        string   `json:"romanized"`
	Keywords         []string `json:"keywords"`
	DetectedLanguage string   `json:"detected_language"`
	SearchIntent     string   `json:"search_intent"`
}

// TranslateSearch translates search queries for multilingual search
func (h *AIHandler) TranslateSearch(c *gin.Context) {
	var req TranslateSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prompt := `Translate the following search query to help with multilingual product search.
Original query: "` + req.Query + `"

Provide translations and search-relevant keywords in JSON format:
{
    "japanese": "Japanese translation/keywords",
    "english": "English translation/keywords",
    "romanized": "Romanized version if applicable",
    "keywords": ["keyword1", "keyword2", "keyword3"],
    "detected_language": "original language code (ja/en/etc)",
    "search_intent": "brief description of what user is looking for"
}

For example:
- If query is "スマホ", return english "smartphone", keywords ["phone", "mobile", "iPhone", "Android"]
- If query is "laptop", return japanese "ノートパソコン", keywords ["PC", "MacBook", "computer"]
- Be creative with synonyms and related terms for better search results`

	// Call Gemini API
	result, err := h.callGeminiAPI(prompt)
	if err != nil {
		// Fallback response
		c.JSON(http.StatusOK, TranslateSearchResponse{
			Japanese:         req.Query,
			English:          req.Query,
			Romanized:        req.Query,
			Keywords:         []string{req.Query},
			DetectedLanguage: "unknown",
			SearchIntent:     req.Query,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AIHandler) callGeminiAPI(prompt string) (TranslateSearchResponse, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-flash-latest:generateContent?key=" + h.geminiAPIKey

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.7,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return TranslateSearchResponse{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return TranslateSearchResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TranslateSearchResponse{}, err
	}

	// Parse Gemini response
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return TranslateSearchResponse{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return TranslateSearchResponse{}, err
	}

	// Extract JSON from response
	text := geminiResp.Candidates[0].Content.Parts[0].Text

	// Remove markdown code blocks if present
	text = cleanJSONResponse(text)

	var result TranslateSearchResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return TranslateSearchResponse{}, err
	}

	return result, nil
}

func cleanJSONResponse(text string) string {
	// Remove ```json and ``` markers
	textBytes := []byte(text)
	textBytes = bytes.TrimPrefix(textBytes, []byte("```json"))
	textBytes = bytes.TrimPrefix(textBytes, []byte("```"))
	textBytes = bytes.TrimSuffix(textBytes, []byte("```"))
	return string(bytes.TrimSpace(textBytes))
}
