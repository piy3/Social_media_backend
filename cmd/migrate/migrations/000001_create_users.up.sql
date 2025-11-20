CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email CITEXT UNIQUE NOT NULL,
  password BYTEA NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);



CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);