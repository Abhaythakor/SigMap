package vulnintel

// CalculateRisk determines the overall risk level based on signals.
func CalculateRisk(p VulnProfile) string {
	if p.CVECount == 0 {
		return "Low"
	}

	// Critical triggers
	if p.ExploitedInWild {
		return "Critical"
	}

	// High triggers
	if p.ExploitAvailable || p.HighSeverityCount > 0 {
		return "High"
	}

	// Medium triggers
	if p.POCAvailable || p.CVECount > 5 {
		return "Medium"
	}

	return "Low"
}
