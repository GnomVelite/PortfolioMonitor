package repositories

import (
	"database/sql"
	//"errors"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"
)

type authProviderRepository struct {
	db *sql.DB
}

func NewAuthProviderRepository(db *sql.DB) AuthProviderRepository {
	return &authProviderRepository{db}
}

func (r *authProviderRepository) CreateAuthProvider(ap *models.AuthProvider) error {
	query := `
        INSERT INTO auth_providers (user_id, provider, provider_id, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
	return r.db.QueryRow(query, ap.UserID, ap.Provider, ap.ProviderID, ap.Email, ap.PasswordHash).
		Scan(&ap.ID)
}

func (r *authProviderRepository) GetAuthProvider(provider, identifier string) (*models.AuthProvider, error) {
	ap := &models.AuthProvider{}
	var query string
	var err error
	if provider == "local" {
		query = `
            SELECT id, user_id, provider, provider_id, email, password_hash, created_at, updated_at
            FROM auth_providers
            WHERE provider = $1 AND email = $2`
		err = r.db.QueryRow(query, provider, identifier).
			Scan(&ap.ID, &ap.UserID, &ap.Provider, &ap.ProviderID, &ap.Email, &ap.PasswordHash, &ap.CreatedAt, &ap.UpdatedAt)
	} else {
		query = `
            SELECT id, user_id, provider, provider_id, email, password_hash, created_at, updated_at
            FROM auth_providers
            WHERE provider = $1 AND provider_id = $2`
		err = r.db.QueryRow(query, provider, identifier).
			Scan(&ap.ID, &ap.UserID, &ap.Provider, &ap.ProviderID, &ap.Email, &ap.PasswordHash, &ap.CreatedAt, &ap.UpdatedAt)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return ap, nil
}

func (r *authProviderRepository) UpdateAuthProvider(ap *models.AuthProvider) error {
	query := `
        UPDATE auth_providers SET email = $1, password_hash = $2, updated_at = NOW()
        WHERE id = $3`
	_, err := r.db.Exec(query, ap.Email, ap.PasswordHash, ap.ID)
	return err
}
