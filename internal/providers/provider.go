package providers

type PaymentProvider interface {
    CreatePayment(req *PaymentRequest) (*PaymentResponse, error)
    CancelPayment(providerRef string) error
    RefundPayment(providerRef string, amount float64) error
    VerifyWebhook(signature, payload string) bool
}

type PaymentRequest struct {
    OrderID       string
    Amount        float64
    Currency      string
    CustomerEmail string
    ReturnURL     string
    Metadata      map[string]interface{}
}

type PaymentResponse struct {
    ProviderRef string
    PaymentURL  string
    Status      string
}