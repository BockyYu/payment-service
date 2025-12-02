package main

import (
	"log"

	"payment-gateway/internal/config"
	"payment-gateway/internal/handlers"
	"payment-gateway/internal/middleware"
	"payment-gateway/internal/models"
	"payment-gateway/internal/repository"
	"payment-gateway/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 連接資料庫
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自動遷移
	if err := db.AutoMigrate(&models.Payment{}, &models.Transaction{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// 初始化 Repository
	paymentRepo := repository.NewPaymentRepository(db)

	// 初始化 Service
	providerService := services.NewProviderService(cfg)
	paymentService := services.NewPaymentService(paymentRepo, providerService)

	// 初始化 Handler
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	webhookHandler := handlers.NewWebhookHandler(paymentService)

	// 設定 Gin
	if cfg.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 支付相關 (需要認證)
		payments := v1.Group("/payments")
		payments.Use(middleware.AuthMiddleware(cfg.App.APIKey))
		{
			payments.POST("", paymentHandler.CreatePayment)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.GET("", paymentHandler.ListPayments)
			payments.POST("/:id/cancel", paymentHandler.CancelPayment)
			payments.POST("/:id/refund", paymentHandler.RefundPayment)
		}

		// Webhook (不需要認證,但需要驗證簽名)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/adyen", webhookHandler.AdyenWebhook)
			webhooks.POST("/stripe", webhookHandler.StripeWebhook)
		}
	}

	// 啟動服務
	log.Printf("🚀 Payment Gateway started on :%s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
