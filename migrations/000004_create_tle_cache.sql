-- 创建TLE缓存表
CREATE TABLE IF NOT EXISTS tle_cache (
    norad_id INT PRIMARY KEY,
    name VARCHAR(100),
    tle_line1 VARCHAR(70),
    tle_line2 VARCHAR(70),
    epoch TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tle_cache_epoch (epoch)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
