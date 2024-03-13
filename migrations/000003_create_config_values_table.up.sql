BEGIN;

CREATE TABLE config_values (
    id SERIAL PRIMARY KEY,
    config_key_id integer REFERENCES config_keys
        ON DELETE CASCADE,
    environment_id integer REFERENCES environments
        ON DELETE CASCADE,

    str_value TEXT,
    int_value INTEGER,
    float_value FLOAT,
    bool_value BOOLEAN,

    inserted_at timestamp DEFAULT current_timestamp,
    updated_at timestamp
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
