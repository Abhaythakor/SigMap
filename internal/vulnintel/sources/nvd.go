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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	popular := map[string]bool{
		"Nginx": true, "Apache": true, "WordPress": true, "PHP": true, "OpenSSL": true,
	}
	
	if !popular[technology] {
		return nil, nil
	}

	bugTypes := []string{"XSS", "SQL Injection", "Remote Code Execution", "Buffer Overflow", "Auth Bypass"}

	findings := []vulnintel.VulnFinding{
		{
			CVE:              fmt.Sprintf("CVE-2023-%d", rng.Intn(9000)+1000),
			Severity:         float64(rng.Intn(4) + 6),
			Description:      fmt.Sprintf("A critical vulnerability in %s allows an attacker to perform %s via a crafted request.", technology, bugTypes[rng.Intn(len(bugTypes))]),
			BugType:          bugTypes[rng.Intn(len(bugTypes))],
			ExploitAvailable: rng.Float32() > 0.7,
			POCAvailable:     rng.Float32() > 0.5,
			PublishedAt:      time.Now().AddDate(0, 0, -rng.Intn(100)),
		},
	}

	return findings, nil
}
