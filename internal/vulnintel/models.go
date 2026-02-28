package vulnintel

import "time"

// VulnProfile represents the security posture of a technology.
type VulnProfile struct {
	Technology        string    `json:"technology"`
	CVECount          int       `json:"cve_count"`
	HighSeverityCount int       `json:"high_severity_count"`
	ExploitAvailable  bool      `json:"exploit_available"`
	POCAvailable      bool      `json:"poc_available"`
	ExploitedInWild   bool      `json:"exploited_in_wild"`
	RiskLevel         string    `json:"risk_level"` // Low, Medium, High, Critical
	LastChecked       time.Time `json:"last_checked"`
}

// VulnFinding represents a single vulnerability result from a source.
type VulnFinding struct {
	CVE               string  `json:"cve"`
	Severity          float64 `json:"severity"`
	ExploitAvailable  bool    `json:"exploit_available"`
	POCAvailable      bool    `json:"poc_available"`
	ExploitedInWild   bool    `json:"exploited_in_wild"`
}

// SourceConnector defines the interface for vulnerability data providers.
type SourceConnector interface {
	GetName() string
	FetchFindings(technology string) ([]VulnFinding, error)
}
