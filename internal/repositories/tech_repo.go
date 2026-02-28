package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TechListItem struct {
	ID               int
	Name             string
	Category         string
	DomainCount      int
	Confidence       int
	RiskLevel        string
	Icon             string
	CVECount         int
	ExploitAvailable bool
}

type TechRepository struct {
	Pool *pgxpool.Pool
}

func NewTechRepository(pool *pgxpool.Pool) *TechRepository {
	return &TechRepository{Pool: pool}
}

func (r *TechRepository) List(ctx context.Context, limit, offset int) ([]TechListItem, error) {
	query := `
		SELECT 
			t.id, t.name, 
			COALESCE(c.name, 'Uncategorized') as category,
			COUNT(DISTINCT det.domain_id) as domain_count,
			COALESCE(AVG(det.confidence), 0)::INT as avg_conf,
			COALESCE(vp.risk_level, t.risk_level) as risk_level,
			COALESCE(t.icon, ''),
			COALESCE(vp.cve_count, 0) as cve_count,
			COALESCE(vp.exploit_available, FALSE) as exploit_available
		FROM technologies t
		LEFT JOIN technology_categories tc ON t.id = tc.technology_id
		LEFT JOIN categories c ON tc.category_id = c.id
		LEFT JOIN detections det ON t.id = det.technology_id
		LEFT JOIN technology_vuln_profile vp ON t.name = vp.technology
		GROUP BY t.id, c.name, vp.risk_level, vp.cve_count, vp.exploit_available
		ORDER BY domain_count DESC, t.name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TechListItem
	for rows.Next() {
		var item TechListItem
		err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.DomainCount, &item.Confidence, &item.RiskLevel, &item.Icon, &item.CVECount, &item.ExploitAvailable)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
