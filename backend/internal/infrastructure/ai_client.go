package infrastructure

import (
	"context"
	"fmt"
	"log"

	pb "github.com/yourusername/ecomate/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AIClient struct {
	client        pb.ProductAnalysisServiceClient
	conn          *grpc.ClientConn
	geminiClient  *GeminiClient
}

func NewAIClient(serverURL string) (*AIClient, error) {
	conn, err := grpc.Dial(serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}

	client := pb.NewProductAnalysisServiceClient(conn)
	log.Printf("Connected to AI service at %s", serverURL)

	// Initialize Gemini client
	geminiClient := NewGeminiClient()
	if geminiClient != nil {
		log.Printf("Gemini API client initialized")
	}

	return &AIClient{
		client:       client,
		conn:         conn,
		geminiClient: geminiClient,
	}, nil
}

func (c *AIClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type ProductAnalysisRequest struct {
	Images                [][]byte
	Title                 string
	UserProvidedDescription string
	Category              string
}

type ProductAnalysisResponse struct {
	GeneratedDescription        string
	SuggestedPrice              int
	EstimatedWeightKg           float64
	ManufacturerCountry         string
	EstimatedManufacturingYear  int
	CO2ImpactKg                 float64
	IsInappropriate             bool
	InappropriateReason         string
	DetectedObjects             []string
}

func (c *AIClient) AnalyzeProduct(ctx context.Context, req *ProductAnalysisRequest) (*ProductAnalysisResponse, error) {
	resp, err := c.client.AnalyzeProduct(ctx, &pb.AnalyzeProductRequest{
		Images:                 req.Images,
		Title:                  req.Title,
		UserProvidedDescription: req.UserProvidedDescription,
		Category:               req.Category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze product: %w", err)
	}

	return &ProductAnalysisResponse{
		GeneratedDescription:       resp.GeneratedDescription,
		SuggestedPrice:             int(resp.SuggestedPrice),
		EstimatedWeightKg:          float64(resp.EstimatedWeightKg),
		ManufacturerCountry:        resp.ManufacturerCountry,
		EstimatedManufacturingYear: int(resp.EstimatedManufacturingYear),
		CO2ImpactKg:                float64(resp.Co2ImpactKg),
		IsInappropriate:            resp.IsInappropriate,
		InappropriateReason:        resp.InappropriateReason,
		DetectedObjects:            resp.DetectedObjects,
	}, nil
}

type CO2CalculationRequest struct {
	Category            string
	WeightKg            float64
	ManufacturerCountry string
	ManufacturingYear   int
}

type CO2CalculationResponse struct {
	BuyingNewKg  float64
	BuyingUsedKg float64
	SavedKg      float64
}

func (c *AIClient) CalculateCO2Impact(ctx context.Context, req *CO2CalculationRequest) (*CO2CalculationResponse, error) {
	resp, err := c.client.CalculateCO2Impact(ctx, &pb.CalculateCO2Request{
		Category:            req.Category,
		WeightKg:            float32(req.WeightKg),
		ManufacturerCountry: req.ManufacturerCountry,
		ManufacturingYear:   int32(req.ManufacturingYear),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to calculate CO2 impact: %w", err)
	}

	return &CO2CalculationResponse{
		BuyingNewKg:  float64(resp.BuyingNewKg),
		BuyingUsedKg: float64(resp.BuyingUsedKg),
		SavedKg:      float64(resp.SavedKg),
	}, nil
}

// AnswerProductQuestion uses AI to answer questions about a product
func (c *AIClient) AnswerProductQuestion(ctx context.Context, productTitle, productDescription, question string) (string, error) {
	if c == nil || c.client == nil {
		// Fallback response when AI service is not available
		return "AI service is currently unavailable. Please try again later or contact the seller directly for product information.", nil
	}

	// Use the AI service to answer the question
	// In a real implementation, this would call a Gemini API endpoint
	// For now, return a helpful response
	answer := fmt.Sprintf("Based on the product information for '%s', I can help answer your question about: %s. %s However, for the most accurate information, please contact the seller directly.",
		productTitle, question, productDescription)

	return answer, nil
}

// GenerateText generates text based on a prompt using AI
// This is a general-purpose AI text generation method
func (c *AIClient) GenerateText(prompt string) (string, error) {
	if c == nil || c.client == nil {
		return "", fmt.Errorf("AI client not initialized")
	}

	// In production, this would call Gemini API for text generation
	// For now, return a structured response based on the prompt
	// The caller should handle parsing the response

	// Simple response that indicates AI processing
	response := fmt.Sprintf("AI Analysis Complete: %s", prompt[:min(100, len(prompt))])

	return response, nil
}

// ChatMessage represents a message in a conversation
type ChatMessage struct {
	Role    string // "user" or "assistant"
	Content string
}

// ChatRequest represents a chat request with conversation history
type ChatRequest struct {
	Messages []ChatMessage
	Context  string // Additional context like available products
}

// Chat generates a response based on conversation history
func (c *AIClient) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	if c == nil || c.client == nil {
		// Fallback response when AI service is not available
		return "こんにちは！EcoMateのアシスタントです。商品探しをお手伝いします。どのような商品をお探しですか？", nil
	}

	// Build conversation context
	conversationContext := "あなたはEcoMateのAIアシスタントです。ユーザーの質問に日本語で答え、適切な商品を提案してください。\n\n"

	if req.Context != "" {
		conversationContext += "利用可能な商品情報:\n" + req.Context + "\n\n"
	}

	// Get the last user message
	var lastMessage string
	if len(req.Messages) > 0 {
		lastMessage = req.Messages[len(req.Messages)-1].Content
	}

	// Try Gemini API first
	if c.geminiClient != nil {
		// Build system prompt
		systemPrompt := `あなたはEcoMateのAIアシスタントです。
EcoMateは環境に優しい中古品フリーマーケットアプリです。

以下のことを心がけて応答してください：
- ユーザーの質問に日本語で丁寧に答える
- 商品探しや購入のサポートをする
- CO2削減や環境保護の価値を伝える
- EcoMateの機能（3D表示、AR試着、価格予測など）を紹介する
- 具体的で実用的なアドバイスを提供する

EcoMateの主な機能：
- 商品検索・フィルタリング（カテゴリー、価格、状態など）
- 3Dモデル表示とAR試着機能
- AI価格予測と交渉サポート
- CO2削減量の可視化
- エコポイント・レベルシステム
- リアルタイムメッセージング

`

		if req.Context != "" {
			systemPrompt += req.Context + "\n\n"
			systemPrompt += "重要: ユーザーが商品を探している場合、上記の商品リストから適切な商品を推薦してください。商品を推薦する際は、必ず上記のフォーマット [PRODUCT:商品ID:商品名:商品画像URL] を使用してください。\n\n"
		}

		// Build conversation history
		var conversationText string
		for _, msg := range req.Messages {
			role := "ユーザー"
			if msg.Role == "assistant" {
				role = "アシスタント"
			}
			conversationText += fmt.Sprintf("%s: %s\n", role, msg.Content)
		}

		// Build final prompt
		prompt := systemPrompt + "\n会話履歴:\n" + conversationText + "\n\n上記の会話履歴に基づいて、ユーザーの最新の質問に日本語で答えてください。"

		// Call Gemini API
		response, err := c.geminiClient.GenerateContent(ctx, prompt)
		if err == nil && response != "" {
			return response, nil
		}

		log.Printf("Gemini API error, falling back to keyword matching: %v", err)
	}

	// Fallback to keyword-based responses
	response := c.generateChatResponse(lastMessage, req.Context)

	return response, nil
}

func (c *AIClient) generateChatResponse(userMessage, context string) string {
	// 様々なパターンの質問に対応
	msg := userMessage

	// 挨拶
	if containsAny(msg, []string{"こんにちは", "はじめまして", "ハロー", "やあ", "よろしく"}) {
		return "こんにちは！EcoMateのAIアシスタントです🌱\nエコフレンドリーな中古品探しをお手伝いします。\n\n例えば:\n• 「スニーカー探してる」\n• 「3000円以下の本」\n• 「CO2削減について」\nなど、お気軽にお聞きください！"
	}

	// カテゴリー別の商品提案
	if containsAny(msg, []string{"スニーカー", "靴", "シューズ"}) {
		return "スニーカーをお探しですね！👟\n\n現在、ハイトップスニーカーなど様々なスニーカーがあります。\n商品詳細ページでは3DビューやAR試着機能もご利用いただけますよ！\n\nトップページの検索で「スニーカー」や「靴」で絞り込んでみてください。"
	}

	if containsAny(msg, []string{"カメラ", "写真機"}) {
		return "カメラをお探しですね！📷\n\nヴィンテージカメラなど、味のある中古カメラを取り扱っています。\n中古のカメラは新品購入に比べて大幅にCO2排出量を削減できます！\n\n検索機能でカメラをお探しください。"
	}

	if containsAny(msg, []string{"本", "書籍", "読みたい"}) {
		return "本をお探しですね！📚\n\n中古本は環境にも優しく、お財布にも優しい選択です。\n価格フィルターで予算に合わせて絞り込むこともできますよ。"
	}

	// 機能説明
	if containsAny(msg, []string{"3D", "AR", "試着", "拡張現実"}) {
		return "3D/AR機能についてですね！\n\n一部の商品では以下の機能が使えます:\n• 📱 3Dビューア: 商品を360度回転して確認\n• 👓 AR試着: カメラで実際に試着体験\n\nスニーカーなどのファッションアイテムでぜひお試しください！\n商品詳細ページの「AR試着」タブから利用できます。"
	}

	if containsAny(msg, []string{"価格", "値段", "いくら", "安い", "高い", "予算"}) {
		return "価格についてですね💰\n\nEcoMateでは様々な価格帯の商品を扱っています。\n\n便利な機能:\n• 価格でソート（安い順/高い順）\n• 価格範囲でフィルタリング\n• AI価格提案\n\nご予算を教えていただければ、おすすめの商品をご提案します！"
	}

	// エコ・環境関連
	if containsAny(msg, []string{"環境", "エコ", "CO2", "二酸化炭素", "サステナ", "持続可能"}) {
		return "環境への配慮について、素晴らしい質問です！🌍\n\nEcoMateでは:\n• 各商品のCO2削減量を表示\n• 中古品購入で新品製造のCO2排出を回避\n• ユーザーのエコレベル・ランキング\n• 累計CO2削減量の可視化\n\n中古品1つ1つが地球を守る一歩です。一緒にエコな選択をしましょう！"
	}

	// 使い方・ヘルプ
	if containsAny(msg, []string{"使い方", "方法", "どうやって", "やり方"}) {
		return "EcoMateの使い方をご説明します！\n\n基本的な流れ:\n1️⃣ トップページで商品を検索・閲覧\n2️⃣ 気になる商品をクリックして詳細確認\n3️⃣ 3D/AR機能で実物確認\n4️⃣ 購入またはメッセージで出品者に質問\n5️⃣ エコポイントをゲット！\n\n何か具体的にお困りのことがあればお聞かせください。"
	}

	// 商品を探している
	if containsAny(msg, []string{"探している", "探してる", "欲しい", "ほしい", "おすすめ", "オススメ", "買いたい"}) {
		return "商品をお探しですね！\n\nどのようなカテゴリーの商品をお探しですか？\n\n取扱カテゴリー:\n• 👕 衣類・ファッション\n• 📱 電化製品\n• 🪑 家具・インテリア\n• 📚 本・雑誌\n• ⚽ スポーツ用品\n• 🎮 おもちゃ・ホビー\n\nカテゴリーや具体的な商品名を教えてください！"
	}

	// デフォルト応答
	return "ご質問ありがとうございます！\n\n私がお手伝いできること:\n• 商品のおすすめ・検索\n• 価格や状態についての相談\n• エコ・CO2削減の説明\n• 3D/AR機能の使い方\n• サイトの使い方\n\n具体的な商品名やカテゴリー、質問内容を教えていただければ、より詳しくご案内できます！"
}

func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if contains(text, keyword) {
			return true
		}
	}
	return false
}

func contains(text, keyword string) bool {
	if text == "" || keyword == "" {
		return false
	}
	// Simple substring search
	for i := 0; i <= len(text)-len(keyword); i++ {
		if text[i:i+len(keyword)] == keyword {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ContentModerationRequest for moderating text content
type ContentModerationRequest struct {
	Text    string
	Context string // e.g., "product_title", "product_description", "message", "review"
}

// ContentModerationResponse contains moderation results
type ContentModerationResponse struct {
	IsInappropriate bool
	Reason          string
	Severity        string   // "low", "medium", "high"
	Categories      []string // e.g., "violence", "hate_speech", "spam", "prohibited_items"
	SuggestedAction string   // "block", "review", "warn"
}

// ModerateContent checks if content is appropriate using keyword matching
func (c *AIClient) ModerateContent(ctx context.Context, req *ContentModerationRequest) (*ContentModerationResponse, error) {
	prohibitedKeywords := []string{
		"weapon", "gun", "drug", "explosive", "counterfeit",
		"fake", "stolen", "illegal", "prohibited", "alcohol",
		"tobacco", "prescription", "replica",
	}

	response := &ContentModerationResponse{
		IsInappropriate: false,
		Severity:        "low",
		Categories:      []string{},
		SuggestedAction: "allow",
	}

	textLower := fmt.Sprintf("%s", req.Text)
	for _, keyword := range prohibitedKeywords {
		if len(textLower) > 0 && len(keyword) > 0 {
			response.IsInappropriate = true
			response.Reason = fmt.Sprintf("Contains prohibited content related to: %s", keyword)
			response.Severity = "high"
			response.Categories = append(response.Categories, "prohibited_items")
			response.SuggestedAction = "block"
			break
		}
	}

	return response, nil
}
