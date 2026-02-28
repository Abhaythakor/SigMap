package vulnintel

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool       *pgxpool.Pool
	Connectors []SourceConnector
	Correlator *Correlator
}

func NewService(pool *pgxpool.Pool, connectors []SourceConnector) *Service {
	return &Service{
		Pool:       pool,
		Connectors: connectors,
		Correlator: NewCorrelator(),
	}
}

// GetTechVulnProfile retrieves or updates the vulnerability profile for a technology.
func (s *Service) GetTechVulnProfile(ctx context.Context, technology string) (VulnProfile, error) {
	// 1. Try to fetch from DB
	var p VulnProfile
	err := s.Pool.QueryRow(ctx, `
		SELECT technology, cve_count, high_count, exploit_available, poc_available, exploited_in_wild, risk_level, last_checked
		FROM technology_vuln_profile
		WHERE technology = $1
	`, technology).Scan(&p.Technology, &p.CVECount, &p.HighSeverityCount, &p.ExploitAvailable, &p.POCAvailable, &p.ExploitedInWild, &p.RiskLevel, &p.LastChecked)

	if err == nil {
		// Found in DB
		return p, nil
	}

	// 2. Not found, fetch from sources
	return s.RefreshTechProfile(ctx, technology)
}

// RefreshTechProfile forced refresh of a technology's vulnerability data.
func (s *Service) RefreshTechProfile(ctx context.Context, technology string) (VulnProfile, error) {
	log.Printf("Refreshing vulnerability profile for: %s", technology)
	
	var allFindings []VulnFinding
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, conn := range s.Connectors {
		wg.Add(1)
		go func(c SourceConnector) {
			defer wg.Done()
			findings, err := c.FetchFindings(technology)
			if err != nil {
				log.Printf("Source %s error for %s: %v", c.GetName(), technology, err)
				return
			}
			mu.Lock()
			allFindings = append(allFindings, findings...)
			mu.Unlock()
		}(conn)
	}

	wg.Wait()

	profile := s.Correlator.Correlate(technology, allFindings)

	// 3. Store/Update in DB
	_, err := s.Pool.Exec(ctx, `
		INSERT INTO technology_vuln_profile (technology, cve_count, high_count, exploit_available, poc_available, exploited_in_wild, risk_level, last_checked)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (technology) DO UPDATE SET
			cve_count = EXCLUDED.cve_count,
			high_count = EXCLUDED.high_count,
			exploit_available = EXCLUDED.exploit_available,
			poc_available = EXCLUDED.poc_available,
			exploited_in_wild = EXCLUDED.exploited_in_wild,
			risk_level = EXCLUDED.risk_level,
			last_checked = EXCLUDED.last_checked
	`, profile.Technology, profile.CVECount, profile.HighSeverityCount, profile.ExploitAvailable, profile.POCAvailable, profile.ExploitedInWild, profile.RiskLevel, profile.LastChecked)

	return profile, err
}
