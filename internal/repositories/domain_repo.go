package repositories

import (
	"context"
	"strings"

	"github.com/Abhaythakor/SigMap/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DomainRepository struct {
	Pool *pgxpool.Pool
}

func NewDomainRepository(pool *pgxpool.Pool) *DomainRepository {
	return &DomainRepository{Pool: pool}
}

// EnsureDomain checks if a domain exists, otherwise creates it.
func (r *DomainRepository) EnsureDomain(ctx context.Context, name string) (int, error) {
	var id int
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO domains (name, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP)
		ON CONFLICT (name) DO UPDATE SET updated_at = EXCLUDED.updated_at
		RETURNING id
	`, name).Scan(&id)
	return id, err
}

// AddDetection adds a technology detection for a domain.
func (r *DomainRepository) AddDetection(ctx context.Context, domainID int, techName string, url string, version string, confidence int, source string) error {
	// If version is empty, check if techName has it (e.g. "Sentry:6.13.2")
	if version == "" && strings.Contains(techName, ":") {
		parts := strings.SplitN(techName, ":", 2)
		techName = parts[0]
		version = parts[1]
	}

	var techID int
	err := r.Pool.QueryRow(ctx, "SELECT id FROM technologies WHERE name = $1", techName).Scan(&techID)
	if err == pgx.ErrNoRows {
		err = r.Pool.QueryRow(ctx, `
			INSERT INTO technologies (name, description, risk_level)
			VALUES ($1, 'Automatically detected technology', 'Low')
			RETURNING id
		`, techName).Scan(&techID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = r.Pool.Exec(ctx, `
		INSERT INTO detections (domain_id, technology_id, url, version, confidence, source, last_seen)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		ON CONFLICT ON CONSTRAINT unique_detection DO UPDATE SET 
			last_seen = EXCLUDED.last_seen,
			confidence = EXCLUDED.confidence,
			url = EXCLUDED.url
	`, domainID, techID, url, version, confidence, source)
	
	return err
}

// ToggleBookmark toggles the is_bookmarked status of a domain.
func (r *DomainRepository) ToggleBookmark(ctx context.Context, id int) (bool, error) {
	var isBookmarked bool
	err := r.Pool.QueryRow(ctx, `
		UPDATE domains 
		SET is_bookmarked = NOT is_bookmarked, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING is_bookmarked
	`, id).Scan(&isBookmarked)
	return isBookmarked, err
}

// ListAlertChannels returns all configured notification channels.
func (r *DomainRepository) ListAlertChannels(ctx context.Context) ([]models.AlertChannel, error) {
	rows, err := r.Pool.Query(ctx, "SELECT id, name, type, url, is_active, created_at FROM alert_channels ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []models.AlertChannel
	for rows.Next() {
		var c models.AlertChannel
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.URL, &c.IsActive, &c.CreatedAt); err == nil {
			channels = append(channels, c)
		}
	}
	return channels, nil
}

// AddAlertChannel adds a new notification channel.
func (r *DomainRepository) AddAlertChannel(ctx context.Context, name, cType, url string) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO alert_channels (name, type, url) VALUES ($1, $2, $3)", name, cType, url)
	return err
}

// DeleteAlertChannel removes a notification channel.
func (r *DomainRepository) DeleteAlertChannel(ctx context.Context, id int) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM alert_channels WHERE id = $1", id)
	return err
}

// CreateNote adds a note to a domain.
func (r *DomainRepository) CreateNote(ctx context.Context, domainID int, content string, author string) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO notes (domain_id, content, author, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
	`, domainID, content, author)
	return err
}

func (r *DomainRepository) GetNoteByID(ctx context.Context, id int) (NoteListItem, error) {
	var item NoteListItem
	err := r.Pool.QueryRow(ctx, `
		SELECT 
			n.id, 
			COALESCE(d.name, t.name) as target,
			CASE WHEN n.domain_id IS NOT NULL THEN 'Domain' ELSE 'Technology' END as type,
			n.content, n.author, n.updated_at
		FROM notes n
		LEFT JOIN domains d ON n.domain_id = d.id
		LEFT JOIN technologies t ON n.technology_id = t.id
		WHERE n.id = $1
	`, id).Scan(&item.ID, &item.Target, &item.Type, &item.Content, &item.Author, &item.UpdatedAt)
	return item, err
}

// UpdateNote updates an existing note.
func (r *DomainRepository) UpdateNote(ctx context.Context, id int, content string) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE notes SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2
	`, content, id)
	return err
}

// DeleteNote removes a note.
func (r *DomainRepository) DeleteNote(ctx context.Context, id int) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM notes WHERE id = $1", id)
	return err
}
