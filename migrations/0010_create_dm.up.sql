-- 0010_create_dm.up.sql

CREATE TABLE dm_conversations (
    id              TEXT PRIMARY KEY,
    uri             TEXT UNIQUE NOT NULL,
    participant_a   TEXT NOT NULL REFERENCES accounts(id),
    participant_b   TEXT NOT NULL REFERENCES accounts(id),
    last_message_at TIMESTAMPTZ,
    unread_count_a  INTEGER NOT NULL DEFAULT 0,
    unread_count_b  INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (participant_a, participant_b),
    CHECK (participant_a < participant_b)
);

CREATE TABLE dm_messages (
    id              TEXT PRIMARY KEY,
    uri             TEXT UNIQUE NOT NULL,
    conversation_id TEXT NOT NULL REFERENCES dm_conversations(id),
    sender_id       TEXT NOT NULL REFERENCES accounts(id),
    content         TEXT NOT NULL,
    is_local        BOOLEAN NOT NULL DEFAULT TRUE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX dm_messages_conversation ON dm_messages (conversation_id, created_at DESC);
CREATE INDEX dm_conversations_participant_a ON dm_conversations (participant_a);
CREATE INDEX dm_conversations_participant_b ON dm_conversations (participant_b);
