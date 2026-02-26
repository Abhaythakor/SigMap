-- 003_dashboard_materialized_view.sql

CREATE MATERIALIZED VIEW IF NOT EXISTS view_dashboard_stats AS
SELECT 
    (SELECT COUNT(*) FROM detections) as total_detections,
    (SELECT COALESCE(AVG(confidence), 0) FROM detections) as avg_confidence,
    (SELECT COUNT(*) FROM technologies WHERE risk_level IN ('High', 'Critical')) as risky_technologies,
    (SELECT COUNT(*) FROM domains WHERE is_bookmarked = TRUE) as bookmarked_domains,
    CURRENT_TIMESTAMP as last_refreshed;

CREATE UNIQUE INDEX IF NOT EXISTS idx_dashboard_stats_refresh ON view_dashboard_stats(last_refreshed);
