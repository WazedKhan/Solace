CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- indexing to make sorting faster
CREATE INDEX idx_users_created_at ON users(created_at DESC);

CREATE INDEX idx_users_name_trgm
ON users USING gin (name gin_trgm_ops);
