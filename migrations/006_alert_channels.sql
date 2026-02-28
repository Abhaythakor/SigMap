-- 006_alert_channels.sql

CREATE TABLE IF NOT EXISTS alert_channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'webhook', 'slack', 'discord'
    url TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS alert_history (
    id SERIAL PRIMARY KEY,
    channel_id INT REFERENCES alert_channels(id) ON DELETE CASCADE,
    domain_id INT REFERENCES domains(id) ON DELETE SET NULL,
    tech_name VARCHAR(255),
    risk_level VARCHAR(50),
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
