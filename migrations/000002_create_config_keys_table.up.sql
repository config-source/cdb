CREATE TABLE config_keys (
    id SERIAL PRIMARY KEY,
    
    name TEXT NOT NULL CONSTRAINT key_name_not_empty CHECK (name <> ''),
    can_propagate BOOLEAN NOT NULL DEFAULT TRUE,
    value_type integer NOT NULL CONSTRAINT value_type_range CHECK (value_type BETWEEN 0 AND 3),
    
    created_at timestamp DEFAULT current_timestamp
);
