BEGIN;

TRUNCATE webhook_definitions_to_environments;
DROP TABLE IF EXISTS webhook_definitions_to_environments;
DROP TABLE IF EXISTS webhook_definitions;

COMMIT;