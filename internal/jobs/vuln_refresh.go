package jobs

import (
	"context"
	"log"

	"github.com/Abhaythakor/SigMap/internal/vulnintel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VulnRefreshJob struct {
	Pool        *pgxpool.Pool
	VulnService *vulnintel.Service
}

func NewVulnRefreshJob(pool *pgxpool.Pool, svc *vulnintel.Service) *VulnRefreshJob {
	return &VulnRefreshJob{Pool: pool, VulnService: svc}
}

// Run refreshes all vulnerability profiles for technologies that have detections.
func (j *VulnRefreshJob) Run(ctx context.Context) error {
	log.Println("Starting vulnerability profile refresh job...")

	// 1. Get all technologies that exist in detections
	rows, err := j.Pool.Query(ctx, "SELECT DISTINCT name FROM technologies t JOIN detections d ON t.id = d.technology_id")
	if err != nil {
		return err
	}
	defer rows.Close()

	var techs []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		techs = append(techs, name)
	}

	// 2. Refresh each one
	for _, t := range techs {
		_, err := j.VulnService.RefreshTechProfile(ctx, t)
		if err != nil {
			log.Printf("Error refreshing %s: %v", t, err)
		}
	}

	log.Println("Vulnerability profile refresh job completed")
	return nil
}
