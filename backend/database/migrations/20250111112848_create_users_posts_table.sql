-- +goose Up
CREATE TABLE users_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    post_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE SET NULL,
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE SET NULL
);
CREATE TRIGGER users_posts_updated_at AFTER UPDATE ON users_posts BEGIN UPDATE users_posts SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id; END;

-- +goose Down
DROP TRIGGER IF EXISTS users_posts_updated_at;
DROP TABLE IF EXISTS users_posts;
