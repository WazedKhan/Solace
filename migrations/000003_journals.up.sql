CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE moods (
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

INSERT INTO moods (id, name)
VALUES
    (gen_random_uuid(), 'Happy'),
    (gen_random_uuid(), 'Sad'),
    (gen_random_uuid(), 'Angry'),
    (gen_random_uuid(), 'Anxious'),
    (gen_random_uuid(), 'Neutral');

CREATE TABLE journals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    mood_id UUID REFERENCES moods(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT,
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_journals_user_id ON journals(user_id);

CREATE INDEX idx_journals_status ON journals(status);

CREATE INDEX idx_journals_user_status ON journals(user_id, status);
