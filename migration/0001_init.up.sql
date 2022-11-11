create table IF NOT EXISTS activities (
    activity_group_id BIGINT unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at BIGINT unsigned NOT NULL,
    updated_at BIGINT unsigned NOT NULL,
    deleted_at BIGINT unsigned NOT NULL default 0
);