package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Abhaythakor/SigMap/internal/repositories"
)

type IngestionService struct {
	Repo *repositories.DomainRepository
}

func NewIngestionService(repo *repositories.DomainRepository) *IngestionService {
	return &IngestionService{Repo: repo}
}

// IngestSampleData populates the database with mock domains and technologies.
func (s *IngestionService) IngestSampleData(ctx context.Context) error {
	domains := []string{
		"stripe.com", "notion.so", "github.com", "google.com", "facebook.com",
		"amazon.com", "apple.com", "netflix.com", "microsoft.com", "twitter.com",
		"openai.com", "slack.com", "zoom.us", "spotify.com", "airbnb.com",
	}

	techs := []string{
		"React", "Vue.js", "Nginx", "Apache", "PHP", "Node.js", "Express",
		"Go", "PostgreSQL", "MySQL", "Redis", "Cloudflare", "Google Analytics",
		"Segment", "Mixpanel", "Tailwind CSS", "Bootstrap", "Next.js", "WordPress",
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, domainName := range domains {
		domainID, err := s.Repo.EnsureDomain(ctx, domainName)
		if err != nil {
			return err
		}

		s.LookupInfrastructure(ctx, domainID, domainName)

		numTechs := r.Intn(4) + 3
		shuffledTechs := make([]string, len(techs))
		copy(shuffledTechs, techs)
		r.Shuffle(len(shuffledTechs), func(i, j int) {
			shuffledTechs[i], shuffledTechs[j] = shuffledTechs[j], shuffledTechs[i]
		})

		for i := 0; i < numTechs; i++ {
			version := ""
			// HIGHER PROBABILITY OF EMPTY VERSION FOR TESTING (80% empty)
			if r.Float32() > 0.8 {
				major := r.Intn(10) + 1
				minor := r.Intn(20)
				patch := r.Intn(10)
				version = fmt.Sprintf("v%d.%d.%d", major, minor, patch)
			}

			techName := shuffledTechs[i]
			// 10% chance of using "Tech:Version" format
			if version != "" && r.Float32() > 0.9 {
				techName = fmt.Sprintf("%s:%s", techName, version)
				version = ""
			}

			err = s.Repo.AddDetection(ctx, domainID, techName, "https://"+domainName, version, r.Intn(40)+60, "Mock Scanner")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LookupInfrastructure simulates discovery of IP, ASN, and Cloud Provider.
func (s *IngestionService) LookupInfrastructure(ctx context.Context, domainID int, domainName string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	ips := []string{"34.212.12.1", "52.4.15.22", "104.16.24.5", "13.248.155.12"}
	clouds := []string{"AWS", "Google Cloud", "Cloudflare", "Azure"}
	asns := []int{16509, 15169, 13335, 8075}
	orgs := []string{"Amazon.com", "Google LLC", "Cloudflare, Inc.", "Microsoft Corp"}

	idx := rng.Intn(len(ips))
	
	s.Repo.Pool.Exec(ctx, `
		UPDATE domains SET 
			ip_address = $1, 
			cloud_provider = $2, 
			asn = $3, 
			asn_org = $4 
		WHERE id = $5
	`, ips[idx], clouds[idx], asns[idx], orgs[idx], domainID)
}
