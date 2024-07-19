CREATE TABLE webhook_definitions (
    id SERIAL PRIMARY KEY,
    template TEXT NOT NULL CONSTRAINT template_not_empty CHECK (template <> ''),
    url TEXT NOT NULL CONSTRAINT url_not_empty CHECK (url <> ''),
    authz_header TEXT,

    created_at timestamp default current_timestamp
);

CREATE TABLE webhook_definitions_to_environments (
    environment_id integer REFERENCES environments ON DELETE CASCADE,
    webhook_definition_id integer REFERENCES webhook_definitions ON DELETE CASCADE
);