-- 002_add_gin_indexes.sql

-- Enable pg_trgm extension for fuzzy search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add GIN index for technology name search
CREATE INDEX IF NOT EXISTS idx_technologies_name_gin ON technologies USING GIN (name gin_trgm_ops);

-- Add GIN index for domain name search
CREATE INDEX IF NOT EXISTS idx_domains_name_gin ON domains USING GIN (name gin_trgm_ops);
