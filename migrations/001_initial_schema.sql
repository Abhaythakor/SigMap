-- 001_initial_schema.sql

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    priority INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Technologies
CREATE TABLE IF NOT EXISTS technologies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    website TEXT,
    icon VARCHAR(255),
    risk_level VARCHAR(50) DEFAULT 'Low', -- Low, Medium, High, Critical
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Join table: Technology <-> Categories
CREATE TABLE IF NOT EXISTS technology_categories (
    technology_id INT REFERENCES technologies(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (technology_id, category_id)
);

-- Domains
CREATE TABLE IF NOT EXISTS domains (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    is_bookmarked BOOLEAN DEFAULT FALSE,
    viewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Detections (The core join table)
CREATE TABLE IF NOT EXISTS detections (
    id SERIAL PRIMARY KEY,
    domain_id INT REFERENCES domains(id) ON DELETE CASCADE,
    technology_id INT REFERENCES technologies(id) ON DELETE CASCADE,
    url TEXT,
    version VARCHAR(100),
    confidence INT CHECK (confidence >= 0 AND confidence <= 100),
    source VARCHAR(255),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Notes
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    domain_id INT REFERENCES domains(id) ON DELETE CASCADE,
    technology_id INT REFERENCES technologies(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    author VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT target_check CHECK (
        (domain_id IS NOT NULL AND technology_id IS NULL) OR
        (domain_id IS NULL AND technology_id IS NOT NULL)
    )
);

-- Bookmarks (Optionally separate or just a flag on domains, but a table allows bookmarking specific tech too)
CREATE TABLE IF NOT EXISTS bookmarks (
    id SERIAL PRIMARY KEY,
    domain_id INT REFERENCES domains(id) ON DELETE CASCADE,
    technology_id INT REFERENCES technologies(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT bookmark_target_check CHECK (
        (domain_id IS NOT NULL AND technology_id IS NULL) OR
        (domain_id IS NULL AND technology_id IS NOT NULL)
    ),
    UNIQUE (domain_id, technology_id)
);

-- Indexes for performance
CREATE INDEX idx_detections_domain ON detections(domain_id);
CREATE INDEX idx_detections_tech ON detections(technology_id);
CREATE INDEX idx_detections_last_seen ON detections(last_seen);
CREATE INDEX idx_domains_bookmarked ON domains(is_bookmarked);
CREATE INDEX idx_technologies_risk ON technologies(risk_level);
