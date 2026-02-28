package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type DomainListItem struct {
	ID           int
	Name         string
	IsBookmarked bool
	Technologies []TechTag
	Categories   []string
	Confidence   int
	LastSeen     string
	HighRisk     int
	MediumRisk   int
}

type TechTag struct {
	Name    string
	Icon    string
	Version string
}

type DomainFilters struct {
	Search       string
	Category     string
	Confidence   string // High, Medium, Low
	IsBookmarked bool
}

func (r *DomainRepository) buildListQuery(filters DomainFilters, startArg int) (string, []interface{}) {
	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argCount := startArg

	if filters.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("d.name ILIKE $%d", argCount))
		args = append(args, "%"+filters.Search+"%")
		argCount++
	}

	if filters.IsBookmarked {
		whereClauses = append(whereClauses, "d.is_bookmarked = TRUE")
	}

	if filters.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM detections det2
			JOIN technology_categories tc2 ON det2.technology_id = tc2.technology_id
			JOIN categories c2 ON tc2.category_id = c2.id
			WHERE det2.domain_id = d.id AND c2.name = $%d
		)`, argCount))
		args = append(args, filters.Category)
		argCount++
	}

	if filters.Confidence != "" {
		minConf := 0
		maxConf := 100
		switch strings.ToLower(filters.Confidence) {
		case "high":
			minConf = 90
		case "medium":
			minConf = 60
			maxConf = 89
		case "low":
			maxConf = 59
		}
		whereClauses = append(whereClauses, fmt.Sprintf(`(
			SELECT AVG(confidence) FROM detections WHERE domain_id = d.id
		) BETWEEN $%d AND $%d`, argCount, argCount+1))
		args = append(args, minConf, maxConf)
		argCount += 2
	}

	where := strings.Join(whereClauses, " AND ")
	return where, args
}

func (r *DomainRepository) List(ctx context.Context, limit, offset int, filters DomainFilters) ([]DomainListItem, error) {
	where, whereArgs := r.buildListQuery(filters, 3)
	
	fullArgs := append([]interface{}{limit, offset}, whereArgs...)
	
	// Aggregating Tech:Icon:Version
	query := fmt.Sprintf(`
		SELECT 
			d.id, d.name, d.is_bookmarked,
			COALESCE(ARRAY_AGG(DISTINCT t.name || ':' || COALESCE(t.icon, '') || ':' || COALESCE(det.version, '')) FILTER (WHERE t.name IS NOT NULL), '{}') as techs,
			COALESCE(ARRAY_AGG(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL), '{}') as cats,
			COALESCE(AVG(det.confidence), 0)::INT as avg_conf,
			MAX(det.last_seen) as last_seen,
			COUNT(DISTINCT CASE WHEN vp.risk_level IN ('High', 'Critical') THEN t.id END) as high_risk,
			COUNT(DISTINCT CASE WHEN vp.risk_level = 'Medium' THEN t.id END) as med_risk
		FROM domains d
		LEFT JOIN detections det ON d.id = det.domain_id
		LEFT JOIN technologies t ON det.technology_id = t.id
		LEFT JOIN technology_categories tc ON t.id = tc.technology_id
		LEFT JOIN categories c ON tc.category_id = c.id
		LEFT JOIN technology_vuln_profile vp ON t.name = vp.technology
		WHERE %s
		GROUP BY d.id
		ORDER BY d.updated_at DESC
		LIMIT $1 OFFSET $2
	`, where)

	rows, err := r.Pool.Query(ctx, query, fullArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DomainListItem
	for rows.Next() {
		var item DomainListItem
		var rawTechs []string
		var lastSeen *time.Time
		err := rows.Scan(&item.ID, &item.Name, &item.IsBookmarked, &rawTechs, &item.Categories, &item.Confidence, &lastSeen, &item.HighRisk, &item.MediumRisk)
		if err != nil {
			return nil, err
		}

		for _, rt := range rawTechs {
			parts := strings.SplitN(rt, ":", 3)
			tag := TechTag{Name: parts[0]}
			if len(parts) > 1 {
				tag.Icon = parts[1]
			}
			if len(parts) > 2 {
				tag.Version = parts[2]
			}
			item.Technologies = append(item.Technologies, tag)
		}

		if lastSeen != nil {
			item.LastSeen = lastSeen.Format("2006-01-02 15:04:05")
		} else {
			item.LastSeen = "Never"
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *DomainRepository) Count(ctx context.Context, filters DomainFilters) (int, error) {
	where, args := r.buildListQuery(filters, 1)
	query := fmt.Sprintf("SELECT COUNT(*) FROM domains d WHERE %s", where)
	
	var count int
	err := r.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}
