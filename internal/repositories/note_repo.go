package repositories

import (
	"context"
	"time"

)

type NoteListItem struct {
	ID        int
	Target    string // Domain or Tech name
	Type      string // "Domain" or "Technology"
	Content   string
	Author    string
	UpdatedAt time.Time
}

func (r *DomainRepository) ListNotes(ctx context.Context) ([]NoteListItem, error) {
	query := `
		SELECT 
			n.id, 
			COALESCE(d.name, t.name) as target,
			CASE WHEN n.domain_id IS NOT NULL THEN 'Domain' ELSE 'Technology' END as type,
			n.content, n.author, n.updated_at
		FROM notes n
		LEFT JOIN domains d ON n.domain_id = d.id
		LEFT JOIN technologies t ON n.technology_id = t.id
		ORDER BY n.updated_at DESC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []NoteListItem
	for rows.Next() {
		var item NoteListItem
		err := rows.Scan(&item.ID, &item.Target, &item.Type, &item.Content, &item.Author, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
