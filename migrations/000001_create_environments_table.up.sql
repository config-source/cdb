CREATE TABLE environments (
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE CONSTRAINT env_name_not_empty CHECK (name <> ''),
    promotes_to_id integer REFERENCES environments,

    created_at timestamp default current_timestamp
);
