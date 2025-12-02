package services

import (
    "payment-gateway/internal/config"
    "payment-gateway/internal/providers"
)

type ProviderService struct {
    providers map[string]providers.PaymentProvider
}

func NewProviderService(cfg *config.Config) *ProviderService {
    return &ProviderService{
        providers: map[string]providers.PaymentProvider{
            "mock": providers.NewMockProvider(),
            // "adyen":  providers.NewAdyenProvider(cfg.Providers.Adyen),
            // "stripe": providers.NewStripeProvider(cfg.Providers.Stripe),
        },
    }
}

func (s *ProviderService) GetProvider(name string) providers.PaymentProvider {
    if provider, ok := s.providers[name]; ok {
        return provider
    }
    return s.providers["mock"] // 預設返回 mock
}