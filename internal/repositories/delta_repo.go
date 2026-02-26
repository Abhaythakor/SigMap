package repositories

import (
	"context"
	"time"

)

type DeltaListItem struct {
	DomainName   string
	TechName     string
	ChangeType   string // "Added", "Removed", "Upgraded"
	Confidence   int
	DetectedAt   time.Time
}

func (r *DomainRepository) ListDelta(ctx context.Context) ([]DeltaListItem, error) {
	// Simple delta: show all detections created in the last 24 hours as "Added"
	query := `
		SELECT 
			d.name, t.name, 'Added' as change_type, det.confidence, det.created_at
		FROM detections det
		JOIN domains d ON det.domain_id = d.id
		JOIN technologies t ON det.technology_id = t.id
		WHERE det.created_at > NOW() - INTERVAL '24 hours'
		ORDER BY det.created_at DESC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DeltaListItem
	for rows.Next() {
		var item DeltaListItem
		err := rows.Scan(&item.DomainName, &item.TechName, &item.ChangeType, &item.Confidence, &item.DetectedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
