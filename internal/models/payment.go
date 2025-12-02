package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusSucceeded  PaymentStatus = "succeeded"
	StatusFailed     PaymentStatus = "failed"
	StatusCancelled  PaymentStatus = "cancelled"
	StatusRefunded   PaymentStatus = "refunded"
)

type Payment struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	OrderID        string         `gorm:"index;not null" json:"order_id"`
	Amount         float64        `gorm:"not null" json:"amount"`
	Currency       string         `gorm:"size:3;not null" json:"currency"`
	Status         PaymentStatus  `gorm:"index;not null" json:"status"`
	Provider       string         `gorm:"index;not null" json:"provider"` // adyen, stripe
	ProviderRef    string         `gorm:"index" json:"provider_ref"`      // 第三方的 payment ID
	IdempotencyKey string         `gorm:"uniqueIndex" json:"-"`
	CustomerEmail  string         `json:"customer_email"`
	PaymentURL     string         `json:"payment_url,omitempty"`
	Metadata       string         `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Transactions   []Transaction  `json:"transactions,omitempty"`
}

type Transaction struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	PaymentID   uint           `gorm:"index;not null" json:"payment_id"`
	Type        string         `gorm:"not null" json:"type"` // authorize, capture, refund
	Amount      float64        `gorm:"not null" json:"amount"`
	Status      string         `gorm:"not null" json:"status"`
	ProviderRef string         `json:"provider_ref"`
	RawResponse string         `gorm:"type:text" json:"-"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
