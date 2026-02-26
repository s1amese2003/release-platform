CREATE TABLE IF NOT EXISTS release_request (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_no VARCHAR(32) NOT NULL,
    app_name VARCHAR(64) NOT NULL,
    app_version VARCHAR(32),
    target_env VARCHAR(16) NOT NULL,
    upgrade_version VARCHAR(16),
    package_path VARCHAR(512) NOT NULL,
    package_md5 VARCHAR(32),
    build_time DATETIME,
    build_jdk VARCHAR(32),
    submitter VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deployed_at DATETIME,
    INDEX idx_app (app_name, target_env),
    INDEX idx_status (status)
);

CREATE TABLE IF NOT EXISTS sql_audit_result (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_id BIGINT NOT NULL,
    sql_file_path VARCHAR(256),
    line_number INT,
    sql_content TEXT,
    sql_type VARCHAR(32),
    risk_level VARCHAR(8) NOT NULL,
    risk_reason VARCHAR(512),
    suggestion VARCHAR(512),
    reviewer_action VARCHAR(16),
    reviewer VARCHAR(64),
    INDEX idx_request (request_id)
);

CREATE TABLE IF NOT EXISTS config_diff_result (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_id BIGINT NOT NULL,
    config_file VARCHAR(128),
    config_key VARCHAR(256),
    package_value TEXT,
    baseline_value TEXT,
    diff_type VARCHAR(16),
    risk_level VARCHAR(8),
    category VARCHAR(32),
    INDEX idx_request (request_id)
);

CREATE TABLE IF NOT EXISTS env_baseline (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    env_code VARCHAR(16) NOT NULL,
    app_name VARCHAR(64) NOT NULL,
    config_key VARCHAR(256) NOT NULL,
    config_value TEXT,
    is_sensitive TINYINT DEFAULT 0,
    source_request BIGINT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_env_app_key (env_code, app_name, config_key)
);

CREATE TABLE IF NOT EXISTS manual_operation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_id BIGINT NOT NULL,
    operation_desc TEXT,
    operation_type VARCHAR(32),
    exec_status VARCHAR(16) DEFAULT 'PENDING',
    executor VARCHAR(64),
    executed_at DATETIME,
    INDEX idx_request (request_id)
);

CREATE TABLE IF NOT EXISTS dependency_change (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    request_id BIGINT NOT NULL,
    group_id VARCHAR(128),
    artifact_id VARCHAR(128),
    old_version VARCHAR(64),
    new_version VARCHAR(64),
    change_type VARCHAR(16),
    cve_ids TEXT,
    risk_level VARCHAR(8),
    INDEX idx_request (request_id)
);
