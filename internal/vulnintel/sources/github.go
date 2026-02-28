package sources

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Abhaythakor/SigMap/internal/vulnintel"
)

type GithubConnector struct{}

func NewGithubConnector() *GithubConnector {
	return &GithubConnector{}
}

func (c *GithubConnector) GetName() string {
	return "GitHub POC"
}

func (c *GithubConnector) FetchFindings(technology string) ([]vulnintel.VulnFinding, error) {
	// In a real world scenario, this would call the GitHub Search API:
	// https://api.github.com/search/repositories?q=CVE-XXXX-XXXX+POC
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 2))
	
	popular := map[string]bool{
		"Nginx": true, "Apache": true, "WordPress": true, "PHP": true, "OpenSSL": true,
	}
	
	if !popular[technology] {
		return nil, nil
	}

	// Simulating finding a POC for a CVE
	findings := []vulnintel.VulnFinding{
		{
			CVE:              fmt.Sprintf("CVE-2023-%d", rng.Intn(1000)+1000),
			Severity:         float64(rng.Intn(2) + 8), // 8.0 - 9.0
			Description:      fmt.Sprintf("Public POC available on GitHub for %s vulnerability.", technology),
			BugType:          "Public Exploit",
			POCAvailable:     true,
			PublishedAt:      time.Now().AddDate(0, 0, -rng.Intn(30)),
		},
	}

	return findings, nil
}
