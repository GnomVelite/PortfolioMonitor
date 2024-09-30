package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GnomVelite/PortfolioMonitor/services/auth-service/internal/config"
	_ "github.com/lib/pq"
)

var Conn *sql.DB

func InitDB(cfg *config.Config) {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	Conn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = Conn.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
}
