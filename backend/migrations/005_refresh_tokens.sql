CREATE TABLE refresh_tokens (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
