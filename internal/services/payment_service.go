package services

import (
    "errors"
    "fmt"

    "payment-gateway/internal/models"
    "payment-gateway/internal/providers"
    "payment-gateway/internal/repository"
)

type PaymentService interface {
    CreatePayment(req *CreatePaymentRequest) (*models.Payment, error)
    GetPayment(id uint) (*models.Payment, error)
    ListPayments(limit, offset int) ([]models.Payment, error)
    CancelPayment(id uint) error
    RefundPayment(id uint, amount float64) error
}

type CreatePaymentRequest struct {
    OrderID        string
    Amount         float64
    Currency       string
    Provider       string
    CustomerEmail  string
    IdempotencyKey string
}

type paymentService struct {
    repo     repository.PaymentRepository
    provider *ProviderService
}

func NewPaymentService(repo repository.PaymentRepository, provider *ProviderService) PaymentService {
    return &paymentService{
        repo:     repo,
        provider: provider,
    }
}

func (s *paymentService) CreatePayment(req *CreatePaymentRequest) (*models.Payment, error) {
    // 冪等性檢查
    if req.IdempotencyKey != "" {
        if existing, err := s.repo.GetByIdempotencyKey(req.IdempotencyKey); err == nil {
            return existing, nil
        }
    }

    // 驗證
    if req.Amount <= 0 {
        return nil, errors.New("amount must be greater than 0")
    }

    // 建立支付記錄
    payment := &models.Payment{
        OrderID:        req.OrderID,
        Amount:         req.Amount,
        Currency:       req.Currency,
        Status:         models.StatusPending,
        Provider:       req.Provider,
        CustomerEmail:  req.CustomerEmail,
        IdempotencyKey: req.IdempotencyKey,
    }

    // 呼叫第三方支付
    provider := s.provider.GetProvider(req.Provider)
    resp, err := provider.CreatePayment(&providers.PaymentRequest{
        OrderID:       req.OrderID,
        Amount:        req.Amount,
        Currency:      req.Currency,
        CustomerEmail: req.CustomerEmail,
    })

    if err != nil {
        payment.Status = models.StatusFailed
        s.repo.Create(payment)
        return nil, fmt.Errorf("failed to create payment with provider: %w", err)
    }

    payment.ProviderRef = resp.ProviderRef
    payment.PaymentURL = resp.PaymentURL
    payment.Status = models.StatusProcessing

    if err := s.repo.Create(payment); err != nil {
        return nil, err
    }

    return payment, nil
}

func (s *paymentService) GetPayment(id uint) (*models.Payment, error) {
    return s.repo.GetByID(id)
}

func (s *paymentService) ListPayments(limit, offset int) ([]models.Payment, error) {
    return s.repo.List(limit, offset)
}

func (s *paymentService) CancelPayment(id uint) error {
    payment, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }

    if payment.Status != models.StatusPending && payment.Status != models.StatusProcessing {
        return errors.New("payment cannot be cancelled")
    }

    provider := s.provider.GetProvider(payment.Provider)
    if err := provider.CancelPayment(payment.ProviderRef); err != nil {
        return err
    }

    payment.Status = models.StatusCancelled
    return s.repo.Update(payment)
}

func (s *paymentService) RefundPayment(id uint, amount float64) error {
    payment, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }

    if payment.Status != models.StatusSucceeded {
        return errors.New("only succeeded payments can be refunded")
    }

    provider := s.provider.GetProvider(payment.Provider)
    if err := provider.RefundPayment(payment.ProviderRef, amount); err != nil {
        return err
    }

    payment.Status = models.StatusRefunded
    return s.repo.Update(payment)
}