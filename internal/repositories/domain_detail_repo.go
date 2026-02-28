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
}

type DetectionHistory struct {
	TechName   string
	Version    string
	DetectedAt time.Time
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

	// 2. Current Stack (enriched with vuln profiles)
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
				d.CurrentStack = append(d.CurrentStack, t)
			}
		}
	}

	// 3. History
	rowsHistory, err := r.Pool.Query(ctx, `
		SELECT t.name, det.version, det.created_at
		FROM detections det
		JOIN technologies t ON det.technology_id = t.id
		WHERE det.domain_id = $1
		ORDER BY det.created_at DESC
		LIMIT 50
	`, id)
	if err == nil {
		defer rowsHistory.Close()
		for rowsHistory.Next() {
			var h DetectionHistory
			if err := rowsHistory.Scan(&h.TechName, &h.Version, &h.DetectedAt); err == nil {
				d.History = append(d.History, h)
			}
		}
	}

	// 4. Subdomains (Discovery)
	// We look for any domains that end with ".root_domain"
	rowsSubs, err := r.Pool.Query(ctx, `
		SELECT name FROM domains 
		WHERE name LIKE '%.' || $1
		AND id != $2
		ORDER BY name ASC
	`, d.Name, d.ID)
	if err == nil {
		defer rowsSubs.Close()
		for rowsSubs.Next() {
			var name string
			if err := rowsSubs.Scan(&name); err == nil {
				d.Subdomains = append(d.Subdomains, name)
			}
		}
	}

	// 5. Notes
	d.Notes, _ = r.ListNotesForDomain(ctx, id)

	return d, nil
}

func (r *DomainRepository) ListNotesForDomain(ctx context.Context, domainID int) ([]NoteListItem, error) {
	query := `
		SELECT 
			n.id, d.name as target, 'Domain' as type,
			n.content, n.author, n.updated_at
		FROM notes n
		JOIN domains d ON n.domain_id = d.id
		WHERE n.domain_id = $1
		ORDER BY n.updated_at DESC
	`
	rows, err := r.Pool.Query(ctx, query, domainID)
	if err != nil {
		return nil, err
	}
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
