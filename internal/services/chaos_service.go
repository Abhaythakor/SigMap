package services

import (
	"context"
	"log"

	"github.com/Abhaythakor/SigMap/internal/integrations/chaos"
	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type ChaosService struct {
	Repo        *repositories.DomainRepository
	ChaosClient *chaos.Client
}

func NewChaosService(repo *repositories.DomainRepository, chaosClient *chaos.Client) *ChaosService {
	return &ChaosService{Repo: repo, ChaosClient: chaosClient}
}

// DiscoverSubdomains fetches and saves subdomains for a root domain.
func (s *ChaosService) DiscoverSubdomains(ctx context.Context, rootDomain string) ([]string, error) {
	subdomains, err := s.ChaosClient.FetchSubdomains(ctx, rootDomain)
	if err != nil {
		return nil, err
	}

	for _, sub := range subdomains {
		_, err := s.Repo.EnsureDomain(ctx, sub)
		if err != nil {
			log.Printf("Chaos: Failed to save subdomain %s: %v", sub, err)
		}
	}

	return subdomains, nil
}
