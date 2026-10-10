CREATE TABLE IF NOT EXISTS workflows (
    id UUID PRIMARY KEY,
    brief TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('queued', 'processing', 'completed', 'failed')),
    result TEXT,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS workflows_created_at_idx ON workflows (created_at DESC);
