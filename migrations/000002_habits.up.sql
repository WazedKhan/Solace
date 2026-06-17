CREATE TABLE habits (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    image_url TEXT,
    current_streak INTEGER NOT NULL DEFAULT 0,
    last_checked_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE habit_checking (
    habit_id UUID NOT NULL REFERENCES habits(id) ON DELETE CASCADE,
    checked_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (habit_id, checked_date)
);


-- indexing
CREATE INDEX idx_habits_user_id ON habits(user_id);
CREATE INDEX idx_habit_checking_habit_id ON habit_checking(habit_id);
