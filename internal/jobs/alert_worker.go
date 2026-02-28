package jobs

import (
	"context"
	"log"

	"github.com/Abhaythakor/SigMap/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertWorker struct {
	Pool         *pgxpool.Pool
	AlertService *services.AlertService
}

func NewAlertWorker(pool *pgxpool.Pool, alertSvc *services.AlertService) *AlertWorker {
	return &AlertWorker{Pool: pool, AlertService: alertSvc}
}

// Run checks for new risky detections in the last hour and dispatches alerts.
func (w *AlertWorker) Run(ctx context.Context) error {
	log.Println("Starting alert worker pass...")

	query := `
		SELECT d.name, t.name, COALESCE(vp.risk_level, t.risk_level) as risk
		FROM detections det
		JOIN domains d ON det.domain_id = d.id
		JOIN technologies t ON det.technology_id = t.id
		LEFT JOIN technology_vuln_profile vp ON t.name = vp.technology
		WHERE det.created_at > NOW() - INTERVAL '1 hour'
		AND (COALESCE(vp.risk_level, t.risk_level) IN ('High', 'Critical'))
		-- Ensure we don't alert twice for the same detection
		AND NOT EXISTS (
			SELECT 1 FROM alert_history ah 
			WHERE ah.domain_id = d.id AND ah.tech_name = t.name 
			AND ah.sent_at > det.created_at
		)
	`

	rows, err := w.Pool.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var domain, tech, risk string
		if err := rows.Scan(&domain, &tech, &risk); err == nil {
			log.Printf("ALERT TRIGGERED: %s detected on %s (Risk: %s)", tech, domain, risk)
			w.AlertService.DispatchAlert(ctx, domain, tech, risk)
		}
	}

	return nil
}
