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
			// Initialize severity label if empty
			if f.SeverityLabel == "" {
				f.SeverityLabel = getSeverityLabel(f.Severity)
			}
			cveMap[f.CVE] = f
			continue
		}

		// Merge logic: take the "worst case"
		if f.Severity > existing.Severity {
			existing.Severity = f.Severity
			existing.SeverityLabel = getSeverityLabel(f.Severity)
		}
		if f.Description != "" && existing.Description == "" {
			existing.Description = f.Description
		}
		if f.BugType != "" && existing.BugType == "" {
			existing.BugType = f.BugType
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
		profile.DetailedVulns = append(profile.DetailedVulns, f)
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

func getSeverityLabel(score float64) string {
	switch {
	case score >= 9.0:
		return "Critical"
	case score >= 7.0:
		return "High"
	case score >= 4.0:
		return "Medium"
	default:
		return "Low"
	}
}
