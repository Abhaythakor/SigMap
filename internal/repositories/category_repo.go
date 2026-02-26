package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryListItem struct {
	ID             int
	Name           string
	TechCount      int
	DomainCount    int
	AvgConfidence  float64
	RiskLevel      string
}

type CategoryRepository struct {
	Pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{Pool: pool}
}

func (r *CategoryRepository) List(ctx context.Context) ([]CategoryListItem, error) {
	query := `
		SELECT 
			c.id, c.name, 
			COUNT(DISTINCT tc.technology_id) as tech_count,
			COUNT(DISTINCT det.domain_id) as domain_count,
			COALESCE(AVG(det.confidence), 0) as avg_conf
		FROM categories c
		LEFT JOIN technology_categories tc ON c.id = tc.category_id
		LEFT JOIN technologies t ON tc.technology_id = t.id
		LEFT JOIN detections det ON t.id = det.technology_id
		GROUP BY c.id
		ORDER BY domain_count DESC, c.name ASC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CategoryListItem
	for rows.Next() {
		var item CategoryListItem
		err := rows.Scan(&item.ID, &item.Name, &item.TechCount, &item.DomainCount, &item.AvgConfidence)
		if err != nil {
			return nil, err
		}
		// Mock risk level based on counts or specific categories for now
		item.RiskLevel = "Low"
		if item.DomainCount > 1000 {
			item.RiskLevel = "Medium"
		}
		items = append(items, item)
	}

	return items, nil
}
