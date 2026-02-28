-- 005_domain_infrastructure_metadata.sql

ALTER TABLE domains 
ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45),
ADD COLUMN IF NOT EXISTS cloud_provider VARCHAR(100),
ADD COLUMN IF NOT EXISTS asn INT,
ADD COLUMN IF NOT EXISTS asn_org VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_domains_ip ON domains(ip_address);
CREATE INDEX IF NOT EXISTS idx_domains_cloud ON domains(cloud_provider);
