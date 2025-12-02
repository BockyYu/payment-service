package providers

import (
    "fmt"
    "time"
)

type MockProvider struct{}

func NewMockProvider() PaymentProvider {
    return &MockProvider{}
}

func (p *MockProvider) CreatePayment(req *PaymentRequest) (*PaymentResponse, error) {
    // 模擬支付建立
    providerRef := fmt.Sprintf("mock_%d", time.Now().Unix())
    
    return &PaymentResponse{
        ProviderRef: providerRef,
        PaymentURL:  fmt.Sprintf("https://mock-payment.com/pay/%s", providerRef),
        Status:      "pending",
    }, nil
}

func (p *MockProvider) CancelPayment(providerRef string) error {
    // 模擬取消
    return nil
}

func (p *MockProvider) RefundPayment(providerRef string, amount float64) error {
    // 模擬退款
    return nil
}

func (p *MockProvider) VerifyWebhook(signature, payload string) bool {
    // 模擬驗證
    return true
}