CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    value DOUBLE PRECISION,
    delta BIGINT,
    CONSTRAINT chk_value_delta CHECK ((value IS NULL) OR (delta IS NULL))
); 
