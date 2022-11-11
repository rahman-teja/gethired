create table todos (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    activity_group_id int NOT NULL,
    title TEXT NOT NULL,
    is_active CHAR(1) NOT NULL,
    priority VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL default NULL,
    is_deleted TINYINT(1) default 0,
    INDEX (is_deleted),
    INDEX (activity_group_id),
    FOREIGN KEY (activity_group_id) REFERENCES activites(id)
);