package vulnintel

import (
	"time"
)

// Correlator merges findings from different sources and computes the final profile.
type Correlator struct{}

func NewCorrelator() *Correlator {
	return &Correlator{}
}

// Correlate merges raw findings into a unified VulnProfile.
func (c *Correlator) Correlate(technology string, findings []VulnFinding) VulnProfile {
	profile := VulnProfile{
		Technology:  technology,
		LastChecked: time.Now(),
		RiskLevel:   "Low",
	}

	if len(findings) == 0 {
		return profile
	}

	cveMap := make(map[string]VulnFinding)
	for _, f := range findings {
		existing, ok := cveMap[f.CVE]
		if !ok {
			cveMap[f.CVE] = f
			continue
		}

		// Merge logic: take the "worst case" for each field
		if f.Severity > existing.Severity {
			existing.Severity = f.Severity
		}
		if f.ExploitAvailable {
			existing.ExploitAvailable = true
		}
		if f.POCAvailable {
			existing.POCAvailable = true
		}
		if f.ExploitedInWild {
			existing.ExploitedInWild = true
		}
		cveMap[f.CVE] = existing
	}

	profile.CVECount = len(cveMap)
	for _, f := range cveMap {
		if f.Severity >= 8.0 {
			profile.HighSeverityCount++
		}
		if f.ExploitAvailable {
			profile.ExploitAvailable = true
		}
		if f.POCAvailable {
			profile.POCAvailable = true
		}
		if f.ExploitedInWild {
			profile.ExploitedInWild = true
		}
	}

	profile.RiskLevel = CalculateRisk(profile)
	return profile
}
