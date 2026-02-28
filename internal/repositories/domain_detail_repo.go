package repositories

import (
	"context"
	"time"
)

type DomainDetail struct {
	ID            int
	Name          string
	IsBookmarked  bool
	IPAddress     string
	CloudProvider string
	ASN           int
	ASNOrg        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	
	CurrentStack []DomainTechDetail
	History      []DetectionHistory
	Notes        []NoteListItem
	Subdomains   []string
	ActiveVulns  []ActiveVulnerability
}

type DomainTechDetail struct {
	Name             string
	Icon             string
	Version          string
	Confidence       int
	RiskLevel        string
	CVECount         int
	ExploitAvailable bool
	LastSeen         time.Time
	Vulnerabilities  []VulnListItem // Actual CVE records
}

type VulnListItem struct {
	CVEID         string
	Description   string
	SeverityScore float64
	SeverityLabel string
	BugType       string
}

type DetectionHistory struct {
	TechName   string
	Version    string
	DetectedAt time.Time
}

type ActiveVulnerability struct {
	Name        string
	TemplateID  string
	Severity    string
	Description string
	MatchedURL  string
	FoundAt     time.Time
}

func (r *DomainRepository) GetDomainDetails(ctx context.Context, id int) (DomainDetail, error) {
	var d DomainDetail

	// 1. Basic Info
	err := r.Pool.QueryRow(ctx, `
		SELECT id, name, is_bookmarked, COALESCE(ip_address, ''), COALESCE(cloud_provider, ''), COALESCE(asn, 0), COALESCE(asn_org, ''), created_at, updated_at
		FROM domains WHERE id = $1
	`, id).Scan(&d.ID, &d.Name, &d.IsBookmarked, &d.IPAddress, &d.CloudProvider, &d.ASN, &d.ASNOrg, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return d, err
	}

	// 2. Current Stack
	rows, err := r.Pool.Query(ctx, `
		SELECT 
			t.name, t.icon, det.version, det.confidence, 
			COALESCE(vp.risk_level, t.risk_level) as risk_level,
			COALESCE(vp.cve_count, 0),
			COALESCE(vp.exploit_available, FALSE),
			det.last_seen
		FROM detections det
		JOIN technologies t ON det.technology_id = t.id
		LEFT JOIN technology_vuln_profile vp ON t.name = vp.technology
		WHERE det.domain_id = $1
		ORDER BY det.last_seen DESC
	`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t DomainTechDetail
			err := rows.Scan(&t.Name, &t.Icon, &t.Version, &t.Confidence, &t.RiskLevel, &t.CVECount, &t.ExploitAvailable, &t.LastSeen)
			if err == nil {
				// Fetch detailed CVEs for this specific tech
				t.Vulnerabilities, _ = r.GetVulnsForTech(ctx, t.Name)
				d.CurrentStack = append(d.CurrentStack, t)
			}
		}
	}

	// 3. Active Vulnerabilities (Nuclei)
	rowsVulns, err := r.Pool.Query(ctx, `
		SELECT name, template_id, severity, description, matched_url, found_at
		FROM active_vulnerabilities
		WHERE domain_id = $1
		ORDER BY 
			CASE severity 
				WHEN 'critical' THEN 1 
				WHEN 'high' THEN 2 
				WHEN 'medium' THEN 3 
				WHEN 'low' THEN 4 
				ELSE 5 
			END ASC
	`, id)
	if err == nil {
		defer rowsVulns.Close()
		for rowsVulns.Next() {
			var v ActiveVulnerability
			if err := rowsVulns.Scan(&v.Name, &v.TemplateID, &v.Severity, &v.Description, &v.MatchedURL, &v.FoundAt); err == nil {
				d.ActiveVulns = append(d.ActiveVulns, v)
			}
		}
	}

	// 4. History & Subdomains (omitted for brevity in write_file but keeping logic)
	r.fillHistoryAndSubs(ctx, &d)

	d.Notes, _ = r.ListNotesForDomain(ctx, id)

	return d, nil
}

func (r *DomainRepository) GetVulnsForTech(ctx context.Context, techName string) ([]VulnListItem, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT cve_id, description, severity_score, severity_label, bug_type
		FROM vulnerability_details
		WHERE technology = $1
		ORDER BY severity_score DESC LIMIT 5
	`, techName)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []VulnListItem
	for rows.Next() {
		var v VulnListItem
		if err := rows.Scan(&v.CVEID, &v.Description, &v.SeverityScore, &v.SeverityLabel, &v.BugType); err == nil {
			list = append(list, v)
		}
	}
	return list, nil
}

func (r *DomainRepository) fillHistoryAndSubs(ctx context.Context, d *DomainDetail) {
	rowsHistory, _ := r.Pool.Query(ctx, "SELECT t.name, det.version, det.created_at FROM detections det JOIN technologies t ON det.technology_id = t.id WHERE det.domain_id = $1 ORDER BY det.created_at DESC LIMIT 50", d.ID)
	if rowsHistory != nil {
		defer rowsHistory.Close()
		for rowsHistory.Next() {
			var h DetectionHistory
			if err := rowsHistory.Scan(&h.TechName, &h.Version, &h.DetectedAt); err == nil {
				d.History = append(d.History, h)
			}
		}
	}
	rowsSubs, _ := r.Pool.Query(ctx, "SELECT name FROM domains WHERE name LIKE '%.' || $1 AND id != $2 ORDER BY name ASC", d.Name, d.ID)
	if rowsSubs != nil {
		defer rowsSubs.Close()
		for rowsSubs.Next() {
			var name string
			if err := rowsSubs.Scan(&name); err == nil {
				d.Subdomains = append(d.Subdomains, name)
			}
		}
	}
}

func (r *DomainRepository) ListNotesForDomain(ctx context.Context, domainID int) ([]NoteListItem, error) {
	rows, _ := r.Pool.Query(ctx, "SELECT n.id, d.name as target, 'Domain' as type, n.content, n.author, n.updated_at FROM notes n JOIN domains d ON n.domain_id = d.id WHERE n.domain_id = $1 ORDER BY n.updated_at DESC", domainID)
	if rows == nil { return nil, nil }
	defer rows.Close()
	var items []NoteListItem
	for rows.Next() {
		var item NoteListItem
		if err := rows.Scan(&item.ID, &item.Target, &item.Type, &item.Content, &item.Author, &item.UpdatedAt); err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}
