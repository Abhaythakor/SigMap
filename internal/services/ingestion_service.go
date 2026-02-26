package services

import (
	"context"
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

		// Assign 3-6 random technologies to each domain
		numTechs := r.Intn(4) + 3
		shuffledTechs := make([]string, len(techs))
		copy(shuffledTechs, techs)
		r.Shuffle(len(shuffledTechs), func(i, j int) {
			shuffledTechs[i], shuffledTechs[j] = shuffledTechs[j], shuffledTechs[i]
		})

		for i := 0; i < numTechs; i++ {
			err = s.Repo.AddDetection(ctx, domainID, shuffledTechs[i], "https://"+domainName, "v1.0.0", r.Intn(40)+60, "Mock Scanner")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
