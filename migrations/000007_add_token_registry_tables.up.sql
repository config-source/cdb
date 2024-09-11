CREATE TABLE revoked_tokens (
    token TEXT NOT NULL UNIQUE CONSTRAINT token_not_empty CHECK (token <> '')
);

CREATE TABLE api_tokens (
    user_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE CONSTRAINT token_not_empty CHECK (token <> ''),
    created_at timestamp default current_timestamp
);
CREATE INDEX ON api_tokens(user_id);
