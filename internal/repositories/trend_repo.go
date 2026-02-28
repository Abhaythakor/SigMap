package repositories

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TrendStat struct {
	Label string
	Value int
	Trend float64
}

type TrendRepo struct {
	Pool *pgxpool.Pool
}

func NewTrendRepo(pool *pgxpool.Pool) *TrendRepo {
	return &TrendRepo{Pool: pool}
}

func (r *TrendRepo) GetTrends(ctx context.Context) ([]TrendStat, error) {
	var stats []TrendStat

	volume, trend, err := r.calculateVolumeTrend(ctx)
	if err == nil {
		stats = append(stats, TrendStat{Label: "New Detections (30d)", Value: volume, Trend: trend})
	}

	risky, rTrend, err := r.calculateRiskTrend(ctx)
	if err == nil {
		stats = append(stats, TrendStat{Label: "Risky Assets Found", Value: risky, Trend: rTrend})
	}

	cats, cTrend, err := r.calculateCategoryTrend(ctx)
	if err == nil {
		stats = append(stats, TrendStat{Label: "Active Categories", Value: cats, Trend: cTrend})
	}

	return stats, nil
}

func (r *TrendRepo) GetVelocityData(ctx context.Context) ([]TrendPoint, error) {
	query := `
		SELECT TO_CHAR(created_at, 'DD MON') as day, COUNT(*)
		FROM detections
		WHERE created_at > NOW() - INTERVAL '30 days'
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

func (r *TrendRepo) calculateVolumeTrend(ctx context.Context) (int, float64, error) {
	var current, previous int
	err := r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM detections WHERE created_at > NOW() - INTERVAL '30 days'").Scan(&current)
	if err != nil {
		return 0, 0, err
	}
	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM detections WHERE created_at BETWEEN NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'").Scan(&previous)
	if err != nil {
		return current, 0, nil
	}

	return current, calculatePercentageChange(current, previous), nil
}

func (r *TrendRepo) calculateRiskTrend(ctx context.Context) (int, float64, error) {
	var current, previous int
	err := r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM detections d JOIN technologies t ON d.technology_id = t.id WHERE t.risk_level IN ('High', 'Critical') AND d.created_at > NOW() - INTERVAL '30 days'").Scan(&current)
	if err != nil {
		return 0, 0, err
	}
	err = r.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM detections d JOIN technologies t ON d.technology_id = t.id WHERE t.risk_level IN ('High', 'Critical') AND d.created_at BETWEEN NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'").Scan(&previous)
	if err != nil {
		return current, 0, nil
	}

	return current, calculatePercentageChange(current, previous), nil
}

func (r *TrendRepo) calculateCategoryTrend(ctx context.Context) (int, float64, error) {
	var current, previous int
	err := r.Pool.QueryRow(ctx, "SELECT COUNT(DISTINCT tc.category_id) FROM detections d JOIN technology_categories tc ON d.technology_id = tc.technology_id WHERE d.created_at > NOW() - INTERVAL '30 days'").Scan(&current)
	if err != nil {
		return 0, 0, err
	}
	err = r.Pool.QueryRow(ctx, "SELECT COUNT(DISTINCT tc.category_id) FROM detections d JOIN technology_categories tc ON d.technology_id = tc.technology_id WHERE d.created_at BETWEEN NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'").Scan(&previous)
	if err != nil {
		return current, 0, nil
	}

	return current, calculatePercentageChange(current, previous), nil
}

func calculatePercentageChange(current, previous int) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0
		}
		return 0.0
	}
	change := float64(current-previous) / float64(previous) * 100.0
	return math.Round(change*10) / 10
}
