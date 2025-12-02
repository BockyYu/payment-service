package repository

import (
    "payment-gateway/internal/models"
    "gorm.io/gorm"
)

type PaymentRepository interface {
    Create(payment *models.Payment) error
    GetByID(id uint) (*models.Payment, error)
    GetByOrderID(orderID string) (*models.Payment, error)
    GetByIdempotencyKey(key string) (*models.Payment, error)
    Update(payment *models.Payment) error
    List(limit, offset int) ([]models.Payment, error)
}

type paymentRepository struct {
    db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
    return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
    return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByID(id uint) (*models.Payment, error) {
    var payment models.Payment
    err := r.db.Preload("Transactions").First(&payment, id).Error
    return &payment, err
}

func (r *paymentRepository) GetByOrderID(orderID string) (*models.Payment, error) {
    var payment models.Payment
    err := r.db.Where("order_id = ?", orderID).First(&payment).Error
    return &payment, err
}

func (r *paymentRepository) GetByIdempotencyKey(key string) (*models.Payment, error) {
    var payment models.Payment
    err := r.db.Where("idempotency_key = ?", key).First(&payment).Error
    return &payment, err
}

func (r *paymentRepository) Update(payment *models.Payment) error {
    return r.db.Save(payment).Error
}

func (r *paymentRepository) List(limit, offset int) ([]models.Payment, error) {
    var payments []models.Payment
    err := r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&payments).Error
    return payments, err
}