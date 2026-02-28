package models

import "time"

type AlertChannel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AlertPayload struct {
	Domain    string `json:"domain"`
	Tech      string `json:"technology"`
	Risk      string `json:"risk_level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
