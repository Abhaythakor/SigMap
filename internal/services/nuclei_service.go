package services

import (
	"context"
	"log"

	"github.com/Abhaythakor/SigMap/internal/integrations/runner"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type NucleiService struct {
	Repo   *repositories.DomainRepository
	Runner *runner.Runner
}

func NewNucleiService(repo *repositories.DomainRepository, r *runner.Runner) *NucleiService {
	return &NucleiService{Repo: repo, Runner: r}
}

// ScanAndStore runs nuclei and persists findings to the database.
func (s *NucleiService) ScanAndStore(ctx context.Context, domainID int, domainName string) error {
	findings, err := s.Runner.RunNuclei(ctx, domainName)
	if err != nil {
		log.Printf("Nuclei: Scan failed for %s: %v", domainName, err)
		return err
	}

	log.Printf("Nuclei: Found %d issues for %s", len(findings), domainName)

	for _, f := range findings {
		_, err := s.Repo.Pool.Exec(ctx, `
			INSERT INTO active_vulnerabilities (domain_id, template_id, name, severity, description, matched_url)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, domainID, f.TemplateID, f.Info.Name, f.Info.Severity, f.Info.Description, f.MatchedURL)
		if err != nil {
			log.Printf("Nuclei: Error saving finding %s: %v", f.Info.Name, err)
		}
	}

	return nil
}
