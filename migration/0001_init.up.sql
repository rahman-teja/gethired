create table activites (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP default NULL,
    is_deleted TINYINT(1) default 0,
    INDEX (is_deleted)
);