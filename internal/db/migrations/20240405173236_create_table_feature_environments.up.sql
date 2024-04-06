CREATE TABLE IF NOT EXISTS feature_environments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    feature_id VARCHAR(255) NOT NULL,
    db_type VARCHAR(255),
    description TEXT,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS resources (
    id BIGSERIAL PRIMARY KEY,
    feature_environment_id BIGINT NOT NULL REFERENCES feature_environments(id),
    app_name VARCHAR(255) NOT NULL,
    link VARCHAR(255),
    is_auto_update BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
