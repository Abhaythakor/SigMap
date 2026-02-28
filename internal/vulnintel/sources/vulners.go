package sources

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Abhaythakor/SigMap/internal/vulnintel"
)

type VulnersConnector struct{}

func NewVulnersConnector() *VulnersConnector {
	return &VulnersConnector{}
}

func (c *VulnersConnector) GetName() string {
	return "Vulners"
}

func (c *VulnersConnector) FetchFindings(technology string) ([]vulnintel.VulnFinding, error) {
	// Simulation logic for Vulners.com API
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 1))
	
	popular := map[string]bool{
		"Nginx": true, "Apache": true, "WordPress": true, "PHP": true, "OpenSSL": true,
	}
	
	if !popular[technology] {
		return nil, nil
	}

	// Vulners often has more exploit data
	findings := []vulnintel.VulnFinding{
		{
			CVE:              fmt.Sprintf("CVE-%d-%d", 2024, rng.Intn(5000)+100),
			Severity:         float64(rng.Intn(3) + 7), // 7.0 - 9.0
			ExploitAvailable: true,
			POCAvailable:     true,
			ExploitedInWild:  rng.Float32() > 0.8,
		},
	}

	return findings, nil
}
