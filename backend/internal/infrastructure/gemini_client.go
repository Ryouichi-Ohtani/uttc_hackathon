package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiClient() *GeminiClient {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil
	}

	return &GeminiClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
	Tools    []GeminiTool    `json:"tools,omitempty"`
}

type GeminiTool struct {
	GoogleSearchRetrieval *GoogleSearchRetrieval `json:"googleSearchRetrieval,omitempty"`
}

type GoogleSearchRetrieval struct {
	DynamicRetrievalConfig *DynamicRetrievalConfig `json:"dynamicRetrievalConfig,omitempty"`
}

type DynamicRetrievalConfig struct {
	Mode             string  `json:"mode"`
	DynamicThreshold float64 `json:"dynamicThreshold"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text       string      `json:"text,omitempty"`
	InlineData *InlineData `json:"inlineData,omitempty"`
}

type InlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // base64 encoded
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return c.generateContentWithSearch(ctx, prompt, false)
}

// GenerateContentWithSearch generates content with optional Google Search grounding
func (c *GeminiClient) GenerateContentWithSearch(ctx context.Context, prompt string) (string, error) {
	return c.generateContentWithSearch(ctx, prompt, true)
}

func (c *GeminiClient) generateContentWithSearch(ctx context.Context, prompt string, enableSearch bool) (string, error) {
	if c == nil {
		return "", fmt.Errorf("Gemini client not initialized")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent?key=%s", c.apiKey)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	// Add Google Search grounding if enabled
	if enableSearch {
		reqBody.Tools = []GeminiTool{
			{
				GoogleSearchRetrieval: &GoogleSearchRetrieval{
					DynamicRetrievalConfig: &DynamicRetrievalConfig{
						Mode:             "MODE_DYNAMIC",
						DynamicThreshold: 0.7,
					},
				},
			},
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// AnalyzeProductImage analyzes a product image and returns description, tags, and category
func (c *GeminiClient) AnalyzeProductImage(ctx context.Context, imageData string, mimeType string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("Gemini client not initialized")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent?key=%s", c.apiKey)

	prompt := `この商品画像を分析して、以下のJSON形式で情報を返してください。前置きなしで、JSONのみを返してください：

{
  "title": "商品名（日本語、簡潔に）",
  "description": "商品の詳細な説明（日本語、200文字程度）",
  "category": "カテゴリー（clothing/electronics/furniture/books/toys/sports のいずれか）",
  "condition": "推定状態（new/like_new/good/fair のいずれか）",
  "price": 推定価格（円、整数）,
  "weight_kg": 推定重量（kg、小数）,
  "brand": "ブランド名（あれば）",
  "model": "モデル名（あれば）",
  "features": ["特徴1", "特徴2", "特徴3"],
  "pricing_rationale": "価格設定の根拠",
  "category_rationale": "カテゴリ選定の根拠"
}

商品の特徴、ブランド、状態などを詳しく記述してください。JSONのみを返し、前置きや説明文は不要です。`

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						InlineData: &InlineData{
							MimeType: mimeType,
							Data:     imageData,
						},
					},
					{
						Text: prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
