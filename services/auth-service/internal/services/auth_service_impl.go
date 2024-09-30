package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/repositories"
)

type authService struct {
	userRepo         repositories.UserRepository
	authProviderRepo repositories.AuthProviderRepository
	jwtSecret        string
}

func NewAuthService(userRepo repositories.UserRepository, authProviderRepo repositories.AuthProviderRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:         userRepo,
		authProviderRepo: authProviderRepo,
		jwtSecret:        jwtSecret,
	}
}

func (s *authService) Register(user *models.User, password string) error {
	// Check if email is already used
	existingAP, _ := s.authProviderRepo.GetAuthProviderByEmail("local", user.Email)
	if existingAP != nil {
		return errors.New("email already registered")
	}

	// Create user
	err := s.userRepo.CreateUser(user)
	if err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create auth provider
	authProvider := &models.AuthProvider{
		UserID:       user.ID,
		Provider:     "local",
		Email:        user.Email,
		PasswordHash: string(hashedPassword),
	}

	return s.authProviderRepo.CreateAuthProvider(authProvider)
}

func (s *authService) Login(email, password string) (string, error) {
	authProvider, err := s.authProviderRepo.GetAuthProviderByEmail("local", email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(authProvider.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	userIDStr := strconv.Itoa(authProvider.UserID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		Subject:   userIDStr,
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) OAuthLogin(provider string, providerID string, email string, name string) (string, error) {
	// Check if auth provider exists
	authProvider, err := s.authProviderRepo.GetAuthProviderByProviderID(provider, providerID)
	if err != nil {
		return "", err
	}

	var user *models.User

	if authProvider == nil {
		// Create new user
		user = &models.User{
			Name:  name, // Use the name provided
			Email: email,
		}
		err = s.userRepo.CreateUser(user)
		if err != nil {
			return "", err
		}

		// Create new auth provider
		authProvider = &models.AuthProvider{
			UserID:     user.ID,
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
		}
		err = s.authProviderRepo.CreateAuthProvider(authProvider)
		if err != nil {
			return "", err
		}
	} else {
		// Fetch the user associated with the auth provider
		user, err = s.userRepo.GetUserByID(authProvider.UserID)
		if err != nil {
			return "", err
		}
	}

	// Generate JWT
	userIDStr := strconv.Itoa(authProvider.UserID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		Subject:   userIDStr,
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}

func (s *authService) ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
