package handlers

import (
	"strconv"

	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/models"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.authService.Register(input.Email, input.Name, input.Password); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Status(201)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	var input struct {
		Provider   string `json:"provider"`
		ProviderID string `json:"provider_id"`
		Email      string `json:"email"`
		Name       string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.authService.OAuthLogin(input.Provider, input.ProviderID, input.Email, input.Name)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user.ID = userID

	if err := h.authService.UpdateUser(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
