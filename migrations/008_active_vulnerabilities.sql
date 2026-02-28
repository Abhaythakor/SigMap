-- 008_active_vulnerabilities.sql

CREATE TABLE IF NOT EXISTS active_vulnerabilities (
    id SERIAL PRIMARY KEY,
    domain_id INT REFERENCES domains(id) ON DELETE CASCADE,
    template_id VARCHAR(255), -- Nuclei template ID
    name TEXT NOT NULL,       -- Vulnerability name
    severity VARCHAR(50),     -- info, low, medium, high, critical
    description TEXT,
    matched_url TEXT,
    remediation TEXT,
    found_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_active_vuln_domain ON active_vulnerabilities(domain_id);
CREATE INDEX idx_active_vuln_severity ON active_vulnerabilities(severity);
