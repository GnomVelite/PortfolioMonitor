package repositories

import "github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	UpdateUser(user *models.User) error
}
