create table IF NOT EXISTS todos (
    id BIGINT unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    activity_group_id BIGINT unsigned NOT NULL,
    title TEXT NOT NULL,
    is_active CHAR(1) NOT NULL,
    priority VARCHAR(100) NOT NULL,
    created_at BIGINT unsigned NOT NULL,
    updated_at BIGINT unsigned NOT NULL,
    deleted_at BIGINT unsigned NOT NULL default 0,
    INDEX (activity_group_id),
    FOREIGN KEY (activity_group_id) REFERENCES activities(activity_group_id)
);