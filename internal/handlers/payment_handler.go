package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "payment-gateway/internal/services"
    "payment-gateway/pkg/response"
)

type PaymentHandler struct {
    service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
    return &PaymentHandler{service: service}
}

// CreatePayment 建立支付
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
    var req struct {
        OrderID       string  `json:"order_id" binding:"required"`
        Amount        float64 `json:"amount" binding:"required,gt=0"`
        Currency      string  `json:"currency" binding:"required,len=3"`
        Provider      string  `json:"provider" binding:"required"`
        CustomerEmail string  `json:"customer_email" binding:"required,email"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    // 取得 Idempotency Key
    idempotencyKey := c.GetHeader("Idempotency-Key")

    payment, err := h.service.CreatePayment(&services.CreatePaymentRequest{
        OrderID:        req.OrderID,
        Amount:         req.Amount,
        Currency:       req.Currency,
        Provider:       req.Provider,
        CustomerEmail:  req.CustomerEmail,
        IdempotencyKey: idempotencyKey,
    })

    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, payment)
}

// GetPayment 查詢支付
func (h *PaymentHandler) GetPayment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "invalid payment id")
        return
    }

    payment, err := h.service.GetPayment(uint(id))
    if err != nil {
        response.Error(c, http.StatusNotFound, "payment not found")
        return
    }

    response.Success(c, payment)
}

// ListPayments 支付列表
func (h *PaymentHandler) ListPayments(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

    payments, err := h.service.ListPayments(limit, offset)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, gin.H{
        "items":  payments,
        "limit":  limit,
        "offset": offset,
    })
}

// CancelPayment 取消支付
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "invalid payment id")
        return
    }

    if err := h.service.CancelPayment(uint(id)); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(c, gin.H{"message": "payment cancelled"})
}

// RefundPayment 退款
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "invalid payment id")
        return
    }

    var req struct {
        Amount float64 `json:"amount" binding:"required,gt=0"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    if err := h.service.RefundPayment(uint(id), req.Amount); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(c, gin.H{"message": "refund processed"})
}