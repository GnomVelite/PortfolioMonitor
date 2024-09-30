package repositories

import "github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"

type AuthProviderRepository interface {
	CreateAuthProvider(authProvider *models.AuthProvider) error
	GetAuthProviderByEmail(provider, email string) (*models.AuthProvider, error)
	GetAuthProviderByProviderID(provider, providerID string) (*models.AuthProvider, error)
	UpdateAuthProvider(authProvider *models.AuthProvider) error
}
