package models

import "time"

// Category represents a Wappalyzer category.
type Category struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Priority  int       `json:"priority" db:"priority"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Technology represents a Wappalyzer technology signature.
type Technology struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Website     string    `json:"website" db:"website"`
	Icon        string    `json:"icon" db:"icon"`
	RiskLevel   string    `json:"risk_level" db:"risk_level"`
	Cats        []int     `json:"cats,omitempty"` // Matches Wappalyzer 'cats' field
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
