package repositories

import (
	"database/sql"

	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (name, email, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, user.Name, user.Email).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserByProvider(provider, providerID string) (*models.User, error) {
	user := &models.User{}
	query := `
        SELECT u.id, u.name, u.email, u.created_at, u.updated_at
        FROM users u
        INNER JOIN auth_providers ap ON ap.user_id = u.id
        WHERE ap.provider = $1 AND ap.provider_id = $2`
	err := r.db.QueryRow(query, provider, providerID).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	query := `
        UPDATE users SET name = $1, email = $2, updated_at = NOW()
        WHERE id = $3`
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	return err
}
