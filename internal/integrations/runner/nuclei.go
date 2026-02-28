package runner

import (
	"context"
	"encoding/json"
	"log"
	"strings"
)

type NucleiFinding struct {
	TemplateID string `json:"template-id"`
	Info       struct {
		Name        string `json:"name"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	} `json:"info"`
	MatchedURL string `json:"matched-at"`
}

// RunNuclei executes nuclei on a target and returns parsed findings.
func (r *Runner) RunNuclei(ctx context.Context, target string) ([]NucleiFinding, error) {
	log.Printf("Nuclei: Scanning %s", target)

	output, err := r.Execute(ctx, "nuclei", "-u", target, "-json", "-silent", "-tags", "tech,exposure,cve")
	if err != nil {
		return nil, err
	}

	var findings []NucleiFinding
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var f NucleiFinding
		if err := json.Unmarshal([]byte(line), &f); err == nil {
			findings = append(findings, f)
		}
	}

	return findings, nil
}
