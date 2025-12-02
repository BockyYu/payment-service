package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "payment-gateway/internal/services"
    "payment-gateway/pkg/response"
)

type WebhookHandler struct {
    service services.PaymentService
}

func NewWebhookHandler(service services.PaymentService) *WebhookHandler {
    return &WebhookHandler{service: service}
}

// AdyenWebhook Adyen webhook
func (h *WebhookHandler) AdyenWebhook(c *gin.Context) {
    // TODO: 驗證簽名
    // TODO: 處理 webhook 邏輯
    
    var webhook struct {
        EventCode   string `json:"eventCode"`
        PSPReference string `json:"pspReference"`
        Success     bool   `json:"success"`
    }

    if err := c.ShouldBindJSON(&webhook); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    // 處理 webhook
    // ...

    response.Success(c, gin.H{"status": "accepted"})
}

// StripeWebhook Stripe webhook
func (h *WebhookHandler) StripeWebhook(c *gin.Context) {
    // TODO: 驗證簽名
    // TODO: 處理 webhook 邏輯

    response.Success(c, gin.H{"status": "accepted"})
}