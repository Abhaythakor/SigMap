package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Abhaythakor/SigMap/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertService struct {
	Pool *pgxpool.Pool
}

func NewAlertService(pool *pgxpool.Pool) *AlertService {
	return &AlertService{Pool: pool}
}

// DispatchAlert sends an alert to all active channels.
func (s *AlertService) DispatchAlert(ctx context.Context, domainName, techName, riskLevel string) error {
	channels, err := s.getActiveChannels(ctx)
	if err != nil {
		return err
	}

	payload := models.AlertPayload{
		Domain:    domainName,
		Tech:      techName,
		Risk:      riskLevel,
		Message:   fmt.Sprintf("Critical Security Alert: %s detected on %s", techName, domainName),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	jsonPayload, _ := json.Marshal(payload)

	for _, ch := range channels {
		go func(c models.AlertChannel) {
			err := s.sendToChannel(c, jsonPayload)
			if err != nil {
				log.Printf("Failed to send alert to %s: %v", c.Name, err)
			} else {
				s.logAlertHistory(context.Background(), c.ID, domainName, techName, riskLevel)
			}
		}(ch)
	}

	return nil
}

func (s *AlertService) getActiveChannels(ctx context.Context) ([]models.AlertChannel, error) {
	rows, err := s.Pool.Query(ctx, "SELECT id, name, type, url FROM alert_channels WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []models.AlertChannel
	for rows.Next() {
		var c models.AlertChannel
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.URL); err == nil {
			channels = append(channels, c)
		}
	}
	return channels, nil
}

func (s *AlertService) sendToChannel(ch models.AlertChannel, payload []byte) error {
	resp, err := http.Post(ch.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *AlertService) logAlertHistory(ctx context.Context, channelID int, domainName, techName, riskLevel string) {
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO alert_history (channel_id, domain_id, tech_name, risk_level)
		VALUES ($1, (SELECT id FROM domains WHERE name = $2), $3, $4)
	`, channelID, domainName, techName, riskLevel)
	if err != nil {
		log.Printf("Error logging alert history: %v", err)
	}
}
