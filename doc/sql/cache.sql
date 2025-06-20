-- 创建UNLOGGED表
CREATE UNLOGGED TABLE cache (
    key TEXT PRIMARY KEY,
    value TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    ttl INTERVAL
);

-- 存储过程：插入数据
CREATE OR REPLACE FUNCTION insert_cache(
    cache_key TEXT,
    cache_value TEXT,
    cache_ttl INTERVAL
) RETURNS VOID AS $$
BEGIN
    INSERT INTO cache (key, value, ttl)
    VALUES (cache_key, cache_value, cache_ttl)
    ON CONFLICT (key) DO UPDATE SET
        value = EXCLUDED.value,
        created_at = NOW(),
        ttl = EXCLUDED.ttl;
END;
$$ LANGUAGE plpgsql;

-- 存储过程：清理过期数据
CREATE OR REPLACE FUNCTION cleanup_cache() RETURNS VOID AS $$
BEGIN
    DELETE FROM cache WHERE created_at + ttl < NOW();
END;
$$ LANGUAGE plpgsql;

-- 定期调用清理过期数据的存储过程
CREATE OR REPLACE FUNCTION schedule_cleanup() RETURNS TRIGGER AS $$
BEGIN
    PERFORM cleanup_cache();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 设置触发器，定期清理过期数据
CREATE TRIGGER cleanup_trigger
AFTER INSERT OR UPDATE ON cache
FOR EACH STATEMENT
EXECUTE FUNCTION schedule_cleanup();
