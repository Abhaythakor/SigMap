package repositories

import (
	"context"
	"time"

)

type BookmarkListItem struct {
	ID         int
	Type       string // Domain or Technology
	Name       string
	Category   string
	Confidence int
	DateAdded  time.Time
}

func (r *DomainRepository) ListBookmarks(ctx context.Context) ([]BookmarkListItem, error) {
	// For now, we only handle bookmarked domains in the main bookmarks view
	query := `
		SELECT 
			d.id, 'Domain' as type, d.name,
			COALESCE((SELECT c.name FROM technology_categories tc 
			          JOIN categories c ON tc.category_id = c.id 
			          JOIN detections det ON tc.technology_id = det.technology_id
			          WHERE det.domain_id = d.id LIMIT 1), 'Uncategorized') as category,
			COALESCE((SELECT AVG(confidence) FROM detections WHERE domain_id = d.id), 0)::INT as avg_conf,
			d.updated_at
		FROM domains d
		WHERE d.is_bookmarked = TRUE
		ORDER BY d.updated_at DESC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []BookmarkListItem
	for rows.Next() {
		var item BookmarkListItem
		err := rows.Scan(&item.ID, &item.Type, &item.Name, &item.Category, &item.Confidence, &item.DateAdded)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
