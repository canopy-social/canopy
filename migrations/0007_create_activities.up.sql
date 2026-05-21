
CREATE TABLE activities (
    id              TEXT PRIMARY KEY,
    uri             TEXT UNIQUE NOT NULL,
    account_id      TEXT REFERENCES accounts(id),
    type            TEXT NOT NULL,
    object_uri      TEXT,
    target_uri      TEXT,
    raw             JSONB NOT NULL,
    is_local        BOOLEAN NOT NULL DEFAULT TRUE,
    delivery_status TEXT NOT NULL DEFAULT 'pending',
    delivery_attempts INTEGER NOT NULL DEFAULT 0,
    last_delivery_error TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX activities_account ON activities (account_id, created_at DESC);
CREATE INDEX activities_delivery ON activities (delivery_status, created_at) WHERE delivery_status != 'delivered';
CREATE INDEX activities_uri ON activities (uri);
CREATE INDEX activities_object_uri ON activities (object_uri);
