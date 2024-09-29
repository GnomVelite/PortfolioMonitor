package models

import (
	"time"
)

// User model corresponds to the `users` table
type User struct {
	ID        int            `json:"id"`         // id SERIAL PRIMARY KEY
	Name      string         `json:"name"`       // name VARCHAR(100)
	Email     string         `json:"email"`      // email VARCHAR(100) UNIQUE
	CreatedAt time.Time      `json:"created_at"` // created_at TIMESTAMP WITH TIME ZONE
	UpdatedAt time.Time      `json:"updated_at"` // updated_at TIMESTAMP WITH TIME ZONE
	Providers []AuthProvider `json:"providers,omitempty"`
}

// AuthProvider model corresponds to the `auth_providers` table
type AuthProvider struct {
	ID           int       `json:"-"`
	UserID       int       `json:"-"`
	Provider     string    `json:"provider"`
	ProviderID   string    `json:"provider_id,omitempty"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"password,omitempty"` // For input only
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
