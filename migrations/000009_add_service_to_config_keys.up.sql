ALTER TABLE config_keys 
ADD COLUMN service_id integer 
REFERENCES services 
ON DELETE CASCADE 
NOT NULL;

DROP INDEX config_key_name;

ALTER TABLE config_keys ADD UNIQUE (service_id, name);

CREATE INDEX ON config_keys(name);
