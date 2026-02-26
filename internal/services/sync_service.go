package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Abhaythakor/SigMap/internal/models"
)

const (
	fingerprintsURL = "https://raw.githubusercontent.com/projectdiscovery/wappalyzergo/refs/heads/main/fingerprints_data.json"
	categoriesURL   = "https://raw.githubusercontent.com/projectdiscovery/wappalyzergo/refs/heads/main/categories_data.json"
)

type SyncService struct {
	Pool *pgxpool.Pool
}

func NewSyncService(pool *pgxpool.Pool) *SyncService {
	return &SyncService{Pool: pool}
}

// Sync downloads and updates the technology metadata.
func (s *SyncService) Sync(ctx context.Context) error {
	log.Println("Starting Wappalyzer metadata sync...")

	// 1. Sync Categories
	categories, err := s.fetchCategories()
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}
	if err := s.saveCategories(ctx, categories); err != nil {
		return fmt.Errorf("failed to save categories: %w", err)
	}

	// 2. Sync Fingerprints
	fingerprints, err := s.fetchFingerprints()
	if err != nil {
		return fmt.Errorf("failed to fetch fingerprints: %w", err)
	}
	if err := s.saveFingerprints(ctx, fingerprints); err != nil {
		return fmt.Errorf("failed to save fingerprints: %w", err)
	}

	log.Println("Wappalyzer metadata sync completed successfully")
	return nil
}

func (s *SyncService) fetchCategories() (map[string]struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}, error) {
	resp, err := http.Get(categoriesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]struct {
		Name     string `json:"name"`
		Priority int    `json:"priority"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *SyncService) fetchFingerprints() (map[string]models.Technology, error) {
	resp, err := http.Get(fingerprintsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Wappalyzer JSON is structured as { "apps": { "name": { ... } } }
	var wrapper struct {
		Apps map[string]models.Technology `json:"apps"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}

	return wrapper.Apps, nil
}

func (s *SyncService) saveCategories(ctx context.Context, categories map[string]struct {
	Name     string `json:"name"`
	Priority int    `json:"priority"`
}) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for idStr, cat := range categories {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Warning: category ID %s is not an integer", idStr)
			continue
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO categories (id, name, priority)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, priority = EXCLUDED.priority
		`, id, cat.Name, cat.Priority)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *SyncService) saveFingerprints(ctx context.Context, apps map[string]models.Technology) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for name, app := range apps {
		var techID int
		err := tx.QueryRow(ctx, `
			INSERT INTO technologies (name, website, icon, description)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (name) DO UPDATE SET website = EXCLUDED.website, icon = EXCLUDED.icon, description = EXCLUDED.description
			RETURNING id
		`, name, app.Website, app.Icon, app.Description).Scan(&techID)
		if err != nil {
			return err
		}

		// Handle category mapping
		for _, catID := range app.Cats {
			_, err := tx.Exec(ctx, `
				INSERT INTO technology_categories (technology_id, category_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, techID, catID)
			if err != nil {
				// If category ID doesn't exist (it should, as we synced them first), log and continue
				log.Printf("Warning: category ID %d for tech %s does not exist", catID, name)
			}
		}
	}

	return tx.Commit(ctx)
}
