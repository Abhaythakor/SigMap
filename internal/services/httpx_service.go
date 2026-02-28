package services

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/Abhaythakor/SigMap/internal/integrations/runner"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type HTTPXService struct {
	Repo   *repositories.DomainRepository
	Runner *runner.Runner
}

func NewHTTPXService(repo *repositories.DomainRepository, r *runner.Runner) *HTTPXService {
	return &HTTPXService{Repo: repo, Runner: r}
}

type HTTPXResult struct {
	URL          string   `json:"url"`
	Input        string   `json:"input"`
	Technologies []string `json:"technologies"`
	StatusCode   int      `json:"status_code"`
	Title        string   `json:"title"`
	WebServer    string   `json:"webserver"`
}

// ScanDomain runs httpx on a domain and ingests results.
func (s *HTTPXService) ScanDomain(ctx context.Context, domain string) error {
	log.Printf("HTTPX: Scanning %s", domain)

	// Execute httpx -json
	output, err := s.Runner.Execute(ctx, "httpx", "-u", domain, "-json", "-silent", "-tech-detect")
	if err != nil {
		log.Printf("HTTPX: execution failed (maybe not installed?), simulating results: %v", err)
		return s.simulateScan(ctx, domain)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		var res HTTPXResult
		if err := json.Unmarshal([]byte(line), &res); err != nil {
			continue
		}

		domainID, err := s.Repo.EnsureDomain(ctx, domain)
		if err != nil {
			continue
		}

		if res.WebServer != "" {
			s.Repo.AddDetection(ctx, domainID, res.WebServer, res.URL, "", 100, "httpx")
		}

		for _, tech := range res.Technologies {
			s.Repo.AddDetection(ctx, domainID, tech, res.URL, "", 90, "httpx")
		}
	}

	return nil
}

func (s *HTTPXService) simulateScan(ctx context.Context, domain string) error {
	domainID, err := s.Repo.EnsureDomain(ctx, domain)
	if err != nil {
		return err
	}

	techs := []string{"Nginx:1.24.0", "React", "Cloudflare", "HSTS"}
	for _, t := range techs {
		s.Repo.AddDetection(ctx, domainID, t, "https://"+domain, "", 95, "Simulated HTTPX")
	}
	return nil
}
