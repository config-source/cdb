CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CONSTRAINT service_name_not_empty CHECK (name <> ''),
    created_at timestamp default current_timestamp
);

ALTER TABLE environments 
ADD COLUMN service_id integer 
REFERENCES services 
ON DELETE CASCADE 
NOT NULL;

ALTER TABLE environments DROP CONSTRAINT environments_name_key;
DROP INDEX environment_name;

ALTER TABLE environments ADD UNIQUE (service_id, name);

CREATE INDEX ON environments(name);
