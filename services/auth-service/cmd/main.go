package main

import (
	//"github.com/gin-contrib/cors"
	"log"

	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/config"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/db"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/services"
	"github.com/gin-gonic/gin"

	//"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/handlers"
	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/repositories"
)

func main() {
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	defer db.Conn.Close()

	userRepo := repositories.NewUserRepository(db.Conn)
	authProviderRepo := repositories.NewAuthProviderRepository(db.Conn)
	authService := services.NewAuthService(userRepo, authProviderRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()
	//router.Use(cors.Default())
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/oauth/login", authHandler.OAuthLogin)
		authRoutes.PUT("/users/:id", authHandler.UpdateUser)
	}

	log.Println("Starting auth-service on port 8000...")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
