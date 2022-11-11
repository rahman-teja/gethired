create table IF NOT EXISTS activities (
    activity_group_id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL default NULL,
    is_deleted TINYINT(1) default 0,
    INDEX (is_deleted)
);