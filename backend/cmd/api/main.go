package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecomate/backend/internal/config"
	"github.com/yourusername/ecomate/backend/internal/infrastructure"
	"github.com/yourusername/ecomate/backend/internal/interfaces"
	"github.com/yourusername/ecomate/backend/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup database
	db, err := infrastructure.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := infrastructure.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed achievements
	if err := infrastructure.SeedAchievements(db); err != nil {
		log.Fatalf("Failed to seed achievements: %v", err)
	}

	// Initialize repositories
	userRepo := infrastructure.NewUserRepository(db)
	productRepo := infrastructure.NewProductRepository(db)
	purchaseRepo := infrastructure.NewPurchaseRepository(db)
	messageRepo := infrastructure.NewMessageRepository(db)
	sustainabilityRepo := infrastructure.NewSustainabilityRepository(db)
	notificationRepo := infrastructure.NewNotificationRepository(db)
	reviewRepo := infrastructure.NewReviewRepository(db)
	offerRepo := infrastructure.NewOfferRepository(db)
	analyticsRepo := infrastructure.NewAnalyticsRepository(db)
	auctionRepo := infrastructure.NewAuctionRepository(db)
	bidRepo := infrastructure.NewBidRepository(db)
	blockchainRepo := infrastructure.NewBlockchainRepository(db)
	nftRepo := infrastructure.NewNFTRepository(db)
	chatHistoryRepo := infrastructure.NewChatHistoryRepository(db)
	co2GoalRepo := infrastructure.NewCO2GoalRepository(db)
	shippingRepo := infrastructure.NewShippingTrackingRepository(db)
	shippingLabelRepo := infrastructure.NewShippingLabelRepository(db)
	aiAgentRepo := infrastructure.NewAIAgentRepository(db)
	autoPurchaseRepo := infrastructure.NewAutoPurchaseWatchRepository(db)
	autoPurchaseLogRepo := infrastructure.NewAutoPurchaseLogRepository(db)

	// Add database indexes for performance
	if err := infrastructure.AddIndexes(db); err != nil {
		log.Printf("Warning: Failed to add indexes: %v", err)
	}

	// Initialize AI client
	aiClient, err := infrastructure.NewAIClient(cfg.AI.ServiceURL)
	if err != nil {
		log.Printf("Warning: AI service not available: %v", err)
		aiClient = nil
	}
	if aiClient != nil {
		defer aiClient.Close()
	}

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	productUseCase := usecase.NewProductUseCase(productRepo, aiClient)
	purchaseUseCase := usecase.NewPurchaseUseCase(purchaseRepo, productRepo, userRepo, shippingLabelRepo)
	messageUseCase := usecase.NewMessageUseCase(messageRepo, productRepo)
	sustainabilityUseCase := usecase.NewSustainabilityUseCase(sustainabilityRepo, userRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo, purchaseRepo)
	recommendationUseCase := usecase.NewRecommendationUseCase(productRepo, purchaseRepo)
	offerUseCase := usecase.NewOfferUseCase(offerRepo, productRepo, aiClient)
	analyticsUseCase := usecase.NewAnalyticsUseCase(analyticsRepo)
	salesPredictionUseCase := usecase.NewSalesPredictionUseCase(productRepo, purchaseRepo, userRepo)
	auctionUseCase := usecase.NewAuctionUseCase(auctionRepo, bidRepo, productRepo)
	voiceSearchUseCase := usecase.NewVoiceSearchUseCase(productRepo)
	blockchainUseCase := usecase.NewBlockchainUseCase(blockchainRepo, nftRepo, purchaseRepo, productRepo)
	chatHistoryUseCase := usecase.NewChatHistoryUseCase(chatHistoryRepo)
	co2GoalUseCase := usecase.NewCO2GoalUseCase(co2GoalRepo)
	shippingUseCase := usecase.NewShippingTrackingUseCase(shippingRepo)

	// Initialize Gemini client
	geminiClient := infrastructure.NewGeminiClient()
	aiAgentUseCase := usecase.NewAIAgentUseCase(aiAgentRepo, productRepo, offerRepo, purchaseRepo, geminiClient)

	// Initialize Auto-Purchase use case
	autoPurchaseUseCase := usecase.NewAutoPurchaseUseCase(
		autoPurchaseRepo,
		autoPurchaseLogRepo,
		productRepo,
		userRepo,
		purchaseRepo,
		shippingLabelRepo,
		notificationRepo,
	)

	// Connect AI Agent to Offer UseCase (for automatic negotiation)
	offerUseCase.SetAIAgentUseCase(aiAgentUseCase)

	// Initialize handlers
	authHandler := interfaces.NewAuthHandler(authUseCase)
	productHandler := interfaces.NewProductHandler(productUseCase, authUseCase, sustainabilityRepo)
	purchaseHandler := interfaces.NewPurchaseHandler(purchaseUseCase)
	messageHandler := interfaces.NewMessageHandler(messageUseCase, authUseCase)
	sustainabilityHandler := interfaces.NewSustainabilityHandler(sustainabilityUseCase)
	notificationHandler := interfaces.NewNotificationHandler(notificationUseCase)
	reviewHandler := interfaces.NewReviewHandler(reviewUseCase)
	recommendationHandler := interfaces.NewRecommendationHandler(recommendationUseCase)
	offerHandler := interfaces.NewOfferHandler(offerUseCase)
	analyticsHandler := interfaces.NewAnalyticsHandler(analyticsUseCase)
	salesPredictionHandler := interfaces.NewSalesPredictionHandler(salesPredictionUseCase)
	chatbotHandler := interfaces.NewChatbotHandler(aiClient, productRepo)
	auctionHandler := interfaces.NewAuctionHandler(auctionUseCase)
	voiceSearchHandler := interfaces.NewVoiceSearchHandler(voiceSearchUseCase)
	blockchainHandler := interfaces.NewBlockchainHandler(blockchainUseCase)
	adminHandler := interfaces.NewAdminHandler(productUseCase, authUseCase)
	uploadHandler := interfaces.NewUploadHandler()
	aiHandler := interfaces.NewAIHandler()
	chatHistoryHandler := interfaces.NewChatHistoryHandler(chatHistoryUseCase)
	co2GoalHandler := interfaces.NewCO2GoalHandler(co2GoalUseCase)
	shippingHandler := interfaces.NewShippingHandler(shippingUseCase)
	aiAgentHandler := interfaces.NewAIAgentHandler(aiAgentUseCase)
	autoPurchaseHandler := interfaces.NewAutoPurchaseHandler(autoPurchaseUseCase)

	// Setup Gin
	gin.SetMode(cfg.Server.GinMode)
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	v1 := router.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", interfaces.AuthMiddleware(authUseCase), authHandler.GetMe)
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", interfaces.AuthMiddleware(authUseCase), productHandler.Create)
			products.PUT("/:id", interfaces.AuthMiddleware(authUseCase), productHandler.Update)
			products.DELETE("/:id", interfaces.AuthMiddleware(authUseCase), productHandler.Delete)
			products.POST("/:id/ask", productHandler.AskQuestion)
			products.POST("/:id/favorite", interfaces.AuthMiddleware(authUseCase), productHandler.AddFavorite)
			products.DELETE("/:id/favorite", interfaces.AuthMiddleware(authUseCase), productHandler.RemoveFavorite)
		}

		// Purchase routes
		purchases := v1.Group("/purchases")
		purchases.Use(interfaces.AuthMiddleware(authUseCase))
		{
			purchases.POST("", purchaseHandler.Create)
			purchases.GET("", purchaseHandler.List)
			purchases.GET("/:id", purchaseHandler.GetByID)
			purchases.PATCH("/:id/complete", purchaseHandler.Complete)
			purchases.GET("/:id/shipping-label", purchaseHandler.GetShippingLabel)
			purchases.POST("/:id/shipping-label", purchaseHandler.GenerateShippingLabel)
		}

		// Messaging routes
		conversations := v1.Group("/conversations")
		conversations.Use(interfaces.AuthMiddleware(authUseCase))
		{
			conversations.GET("", messageHandler.ListConversations)
			conversations.POST("", messageHandler.CreateConversation)
			conversations.GET("/product/:productId/seller/:sellerId", messageHandler.GetOrCreateConversation)
			conversations.GET("/:id/messages", messageHandler.GetMessages)
			conversations.POST("/:id/messages", messageHandler.SendMessage)
		}

		// WebSocket route
		v1.GET("/ws/conversations/:id", messageHandler.WebSocketHandler)

		// Sustainability routes
		sustainability := v1.Group("/sustainability")
		{
			// Public route - no auth required
			sustainability.GET("/leaderboard", sustainabilityHandler.GetLeaderboard)

			// Protected routes - auth required
			sustainability.GET("/dashboard", interfaces.AuthMiddleware(authUseCase), sustainabilityHandler.GetDashboard)
			sustainability.GET("/favorites", interfaces.AuthMiddleware(authUseCase), sustainabilityHandler.GetUserFavorites)
		}

		// Notification routes
		notifications := v1.Group("/notifications")
		notifications.Use(interfaces.AuthMiddleware(authUseCase))
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
			notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
		}

		// Review routes
		reviews := v1.Group("/reviews")
		{
			reviews.POST("", interfaces.AuthMiddleware(authUseCase), reviewHandler.CreateReview)
			reviews.GET("/products/:id", reviewHandler.GetProductReviews)
		}

		// Recommendation routes
		recommendations := v1.Group("/recommendations")
		recommendations.Use(interfaces.AuthMiddleware(authUseCase))
		{
			recommendations.GET("", recommendationHandler.GetRecommendations)
		}

		// Offer routes
		offers := v1.Group("/offers")
		offers.Use(interfaces.AuthMiddleware(authUseCase))
		{
			offers.POST("", offerHandler.CreateOffer)
			offers.GET("/my", offerHandler.GetMyOffers)
			offers.GET("/products/:id", offerHandler.GetProductOffers)
			offers.GET("/products/:id/ai-suggestion", offerHandler.GetNegotiationSuggestion)
			offers.POST("/:id/ai-negotiate", offerHandler.StartAINegotiation)
			offers.POST("/:id/ai-renegotiate", offerHandler.RetryAINegotiationWithPrompt)
			offers.GET("/:id/market-analysis", offerHandler.GetMarketPriceAnalysis)
			offers.PATCH("/:id/respond", offerHandler.RespondOffer)
		}

		// Chatbot routes (no auth required for public access)
		chatbot := v1.Group("/chatbot")
		{
			chatbot.POST("/chat", chatbotHandler.Chat)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		analytics.Use(interfaces.AuthMiddleware(authUseCase))
		{
			analytics.GET("/user/behavior", analyticsHandler.GetUserBehavior)
			analytics.GET("/popular-products", analyticsHandler.GetPopularProducts)
			analytics.GET("/search-trends", analyticsHandler.GetSearchTrends)
		}

		// Sales Prediction routes
		predictions := v1.Group("/predictions")
		predictions.Use(interfaces.AuthMiddleware(authUseCase))
		{
			predictions.GET("/products/:id/price", salesPredictionHandler.PredictProductPrice)
			predictions.GET("/revenue", salesPredictionHandler.PredictSellerRevenue)
			predictions.GET("/market-trends", salesPredictionHandler.GetMarketTrends)
		}

		// Auction routes
		auctions := v1.Group("/auctions")
		{
			auctions.GET("", auctionHandler.GetActiveAuctions)
			auctions.GET("/:id", auctionHandler.GetAuction)
			auctions.POST("", interfaces.AuthMiddleware(authUseCase), auctionHandler.CreateAuction)
			auctions.POST("/:id/bid", interfaces.AuthMiddleware(authUseCase), auctionHandler.PlaceBid)
			auctions.GET("/:id/bids", auctionHandler.GetAuctionBids)
		}

		// WebSocket route for real-time bidding
		v1.GET("/ws/auctions/:id", auctionHandler.WebSocketHandler)

		// Voice Search routes
		v1.POST("/voice-search/text", voiceSearchHandler.SearchByText)
		v1.POST("/voice-search/audio", interfaces.AuthMiddleware(authUseCase), voiceSearchHandler.SearchByAudio)

		// Blockchain routes
		blockchain := v1.Group("/blockchain")
		blockchain.Use(interfaces.AuthMiddleware(authUseCase))
		{
			blockchain.POST("/record-purchase", blockchainHandler.RecordPurchase)
			blockchain.POST("/mint-nft", blockchainHandler.MintNFT)
			blockchain.GET("/nfts/my", blockchainHandler.GetMyNFTs)
			blockchain.GET("/purchases/:id/transaction", blockchainHandler.GetPurchaseTransaction)
		}

		// Upload routes
		upload := v1.Group("/upload")
		upload.Use(interfaces.AuthMiddleware(authUseCase))
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/images", uploadHandler.UploadMultipleImages)
		}

		// AI routes
		ai := v1.Group("/ai")
		{
			ai.POST("/translate-search", aiHandler.TranslateSearch)
		}

		// AI Agent routes (New!)
		aiAgentHandler.RegisterRoutes(v1, interfaces.AuthMiddleware(authUseCase))

		// Chat history routes
		chatHistory := v1.Group("/chat-history")
		chatHistory.Use(interfaces.AuthMiddleware(authUseCase))
		{
			chatHistory.GET("", chatHistoryHandler.GetHistory)
			chatHistory.DELETE("/:id", chatHistoryHandler.DeleteHistory)
		}

		// CO2 Goal routes
		co2Goals := v1.Group("/co2-goals")
		co2Goals.Use(interfaces.AuthMiddleware(authUseCase))
		{
			co2Goals.POST("", co2GoalHandler.CreateGoal)
			co2Goals.GET("", co2GoalHandler.GetGoal)
		}

		// Shipping routes
		shipping := v1.Group("/shipping")
		shipping.Use(interfaces.AuthMiddleware(authUseCase))
		{
			shipping.POST("", shippingHandler.CreateShipping)
			shipping.GET("/purchase/:purchase_id", shippingHandler.GetShipping)
			shipping.PATCH("/:id/status", shippingHandler.UpdateStatus)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(interfaces.AuthMiddleware(authUseCase))
		admin.Use(interfaces.AdminMiddleware())
		{
			admin.GET("/products", adminHandler.GetAllProducts)
			admin.PUT("/products/:id", adminHandler.AdminUpdateProduct)
			admin.DELETE("/products/:id", adminHandler.AdminDeleteProduct)
			admin.GET("/users", adminHandler.GetAllUsers)
		}

		// Auto-Purchase routes
		autoPurchases := v1.Group("/auto-purchases")
		autoPurchases.Use(interfaces.AuthMiddleware(authUseCase))
		{
			autoPurchases.POST("/authorize-payment", autoPurchaseHandler.AuthorizePayment)
			autoPurchases.POST("", autoPurchaseHandler.CreateWatch)
			autoPurchases.GET("", autoPurchaseHandler.GetUserWatches)
			autoPurchases.GET("/:id", autoPurchaseHandler.GetWatch)
			autoPurchases.DELETE("/:id", autoPurchaseHandler.CancelWatch)
		}

		// Background job endpoint (should be protected in production)
		v1.POST("/auto-purchases/check-and-execute", autoPurchaseHandler.CheckAndExecute)
	}

	// Serve uploaded files
	router.GET("/uploads/:filename", uploadHandler.ServeUploadedFile)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Server.Port)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = ctx // Graceful shutdown logic here

	log.Println("Server exited")
}
