package services

import (
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"
	"github.com/dgrijalva/jwt-go"
)

type AuthService interface {
	Register(email, name, password string) error
	Login(email, password string) (string, error)
	OAuthLogin(provider, providerID, email, name string) (string, error)
	UpdateUser(user *models.User) error
	ValidateToken(tokenString string) (*jwt.StandardClaims, error)
}
