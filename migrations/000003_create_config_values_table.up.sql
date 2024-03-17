BEGIN;

CREATE TABLE config_values (
    id SERIAL PRIMARY KEY,
    config_key_id integer REFERENCES config_keys
        ON DELETE CASCADE
        NOT NULL,
    environment_id integer REFERENCES environments
        ON DELETE CASCADE
        NOT NULL,

    str_value TEXT,
    int_value INTEGER,
    float_value FLOAT,
    bool_value BOOLEAN,

    created_at timestamp DEFAULT current_timestamp,

    UNIQUE (environment_id, config_key_id)
);

CREATE OR REPLACE FUNCTION update_config_value_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = now(); 
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language 'plpgsql';

COMMIT;
