package vulnintel

import "time"

// VulnProfile represents the security posture summary of a technology.
type VulnProfile struct {
	Technology        string    `json:"technology"`
	CVECount          int       `json:"cve_count"`
	HighSeverityCount int       `json:"high_severity_count"`
	ExploitAvailable  bool      `json:"exploit_available"`
	POCAvailable      bool      `json:"poc_available"`
	ExploitedInWild   bool      `json:"exploited_in_wild"`
	RiskLevel         string    `json:"risk_level"`
	LastChecked       time.Time `json:"last_checked"`
	DetailedVulns     []VulnFinding `json:"vulnerabilities,omitempty"`
}

// VulnFinding represents a single detailed vulnerability record.
type VulnFinding struct {
	CVE               string    `json:"cve"`
	Severity          float64   `json:"severity"`
	SeverityLabel     string    `json:"severity_label"`
	Description       string    `json:"description"`
	BugType           string    `json:"bug_type"`
	ExploitAvailable  bool      `json:"exploit_available"`
	POCAvailable      bool      `json:"poc_available"`
	ExploitedInWild   bool      `json:"exploited_in_wild"`
	PublishedAt       time.Time `json:"published_at"`
}

// SourceConnector defines the interface for vulnerability data providers.
type SourceConnector interface {
	GetName() string
	FetchFindings(technology string) ([]VulnFinding, error)
}
