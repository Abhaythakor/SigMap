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
	var p VulnProfile
	err := s.Pool.QueryRow(ctx, `
		SELECT technology, cve_count, high_count, exploit_available, poc_available, exploited_in_wild, risk_level, last_checked
		FROM technology_vuln_profile
		WHERE technology = $1
	`, technology).Scan(&p.Technology, &p.CVECount, &p.HighSeverityCount, &p.ExploitAvailable, &p.POCAvailable, &p.ExploitedInWild, &p.RiskLevel, &p.LastChecked)

	if err == nil {
		// Fetch detailed vulnerabilities
		p.DetailedVulns, _ = s.getDetailedVulns(ctx, technology)
		return p, nil
	}

	return s.RefreshTechProfile(ctx, technology)
}

func (s *Service) getDetailedVulns(ctx context.Context, technology string) ([]VulnFinding, error) {
	rows, err := s.Pool.Query(ctx, `
		SELECT cve_id, description, severity_score, severity_label, bug_type, exploit_available, published_at
		FROM vulnerability_details
		WHERE technology = $1
		ORDER BY severity_score DESC NULLS LAST
	`, technology)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var findings []VulnFinding
	for rows.Next() {
		var f VulnFinding
		err := rows.Scan(&f.CVE, &f.Description, &f.Severity, &f.SeverityLabel, &f.BugType, &f.ExploitAvailable, &f.PublishedAt)
		if err == nil {
			findings = append(findings, f)
		}
	}
	return findings, nil
}

// RefreshTechProfile forced refresh of a technology's vulnerability data.
func (s *Service) RefreshTechProfile(ctx context.Context, technology string) (VulnProfile, error) {
	log.Printf("Refreshing real vulnerability details for: %s", technology)
	
	var allFindings []VulnFinding
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, conn := range s.Connectors {
		wg.Add(1)
		go func(c SourceConnector) {
			defer wg.Done()
			findings, err := c.FetchFindings(technology)
			if err != nil {
				return
			}
			mu.Lock()
			allFindings = append(allFindings, findings...)
			mu.Unlock()
		}(conn)
	}

	wg.Wait()

	profile := s.Correlator.Correlate(technology, allFindings)

	// Update Summary
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

	// Update Details
	for _, v := range profile.DetailedVulns {
		_, _ = s.Pool.Exec(ctx, `
			INSERT INTO vulnerability_details (cve_id, technology, description, severity_score, severity_label, bug_type, exploit_available, published_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (cve_id, technology) DO UPDATE SET
				description = EXCLUDED.description,
				severity_score = EXCLUDED.severity_score,
				severity_label = EXCLUDED.severity_label,
				bug_type = EXCLUDED.bug_type,
				exploit_available = EXCLUDED.exploit_available
		`, v.CVE, technology, v.Description, v.Severity, v.SeverityLabel, v.BugType, v.ExploitAvailable, v.PublishedAt)
	}

	return profile, err
}
