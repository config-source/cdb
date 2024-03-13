BEGIN;

TRUNCATE config_values;
DROP TABLE IF EXISTS config_values;

DROP FUNCTION IF EXISTS update_config_value_updated_at;

COMMIT;
