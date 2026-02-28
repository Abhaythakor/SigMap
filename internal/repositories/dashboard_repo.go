package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardStats struct {
	TotalDetections   int
	HighConfidence    float64
	RiskyTechnologies int
	BookmarkedDomains int
	CriticalTechs     int
}

type TrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type DistributionPoint struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}

type DashboardRepository struct {
	Pool *pgxpool.Pool
}

func NewDashboardRepository(pool *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{Pool: pool}
}

func (r *DashboardRepository) GetStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats

	err := r.Pool.QueryRow(ctx, `
		SELECT total_detections, avg_confidence, risky_technologies, bookmarked_domains 
		FROM view_dashboard_stats
	`).Scan(&stats.TotalDetections, &stats.HighConfidence, &stats.RiskyTechnologies, &stats.BookmarkedDomains)
	
	// Fallback or additional metrics not yet in materialized view
	if err != nil {
		return r.GetStatsRealtime(ctx)
	}

	// Fetch CriticalTechs realtime for now
	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM technology_vuln_profile WHERE risk_level = 'Critical'").Scan(&stats.CriticalTechs)

	return stats, nil
}

func (r *DashboardRepository) GetStatsRealtime(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats

	err := r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM detections").Scan(&stats.TotalDetections)
	if err != nil {
		return stats, err
	}

	err = r.Pool.QueryRow(ctx, "SELECT COALESCE(AVG(confidence), 0) FROM detections").Scan(&stats.HighConfidence)
	if err != nil {
		return stats, err
	}

	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM technology_vuln_profile WHERE risk_level IN ('High', 'Critical')").Scan(&stats.RiskyTechnologies)
	if err != nil {
		return stats, err
	}

	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM domains WHERE is_bookmarked = TRUE").Scan(&stats.BookmarkedDomains)
	if err != nil {
		return stats, err
	}

	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM technology_vuln_profile WHERE risk_level = 'Critical'").Scan(&stats.CriticalTechs)

	return stats, nil
}

func (r *DashboardRepository) RefreshStats(ctx context.Context) error {
	_, err := r.Pool.Exec(ctx, "REFRESH MATERIALIZED VIEW view_dashboard_stats")
	return err
}

func (r *DashboardRepository) GetTrendData(ctx context.Context) ([]TrendPoint, error) {
	query := `
		SELECT TO_CHAR(created_at, 'DD MON') as day, COUNT(*)
		FROM detections
		WHERE created_at > NOW() - INTERVAL '7 days'
		GROUP BY day, DATE_TRUNC('day', created_at)
		ORDER BY DATE_TRUNC('day', created_at) ASC
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []TrendPoint
	for rows.Next() {
		var p TrendPoint
		err := rows.Scan(&p.Date, &p.Count)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}

func (r *DashboardRepository) GetDistributionData(ctx context.Context) ([]DistributionPoint, error) {
	query := `
		SELECT t.name, COUNT(*) * 100.0 / NULLIF((SELECT COUNT(*) FROM detections), 0) as pct
		FROM detections det
		JOIN technologies t ON det.technology_id = t.id
		GROUP BY t.name
		ORDER BY pct DESC
		LIMIT 5
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []DistributionPoint
	for rows.Next() {
		var p DistributionPoint
		err := rows.Scan(&p.Name, &p.Percentage)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	return points, nil
}
