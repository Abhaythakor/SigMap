package sources

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Abhaythakor/SigMap/internal/vulnintel"
)

type NVDConnector struct{}

func NewNVDConnector() *NVDConnector {
	return &NVDConnector{}
}

func (c *NVDConnector) GetName() string {
	return "NVD"
}

func (c *NVDConnector) FetchFindings(technology string) ([]vulnintel.VulnFinding, error) {
	// In a real implementation, this would call:
	// https://services.nvd.nist.gov/rest/json/cves/2.0?keywordSearch=technology
	
	// Simulating findings for MVP/Testing
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Only return findings for known popular techs to simulate realism
	popular := map[string]bool{
		"Nginx": true, "Apache": true, "WordPress": true, "PHP": true, "OpenSSL": true,
	}
	
	if !popular[technology] {
		return nil, nil
	}

	findings := []vulnintel.VulnFinding{
		{
			CVE:              fmt.Sprintf("CVE-%d-%d", 2023, rng.Intn(9000)+1000),
			Severity:         float64(rng.Intn(4) + 6), // 6.0 - 9.0
			ExploitAvailable: rng.Float32() > 0.7,
			POCAvailable:     rng.Float32() > 0.5,
		},
	}

	return findings, nil
}
