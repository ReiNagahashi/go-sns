-- +goose Up
CREATE TABLE posts_backup AS SELECT * FROM posts;

CREATE TABLE posts_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    submitted_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (submitted_by) REFERENCES user(id) ON DELETE SET NULL
);

DROP TABLE posts;

ALTER TABLE posts_new RENAME TO posts;

-- +goose Down
DROP TABLE posts;

ALTER TABLE posts_backup RENAME TO posts;
