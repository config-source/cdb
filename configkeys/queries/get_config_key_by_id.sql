SELECT 
    config_keys.*,
    services.name as service_name 
FROM config_keys
INNER JOIN services ON services.id = config_keys.service_id
WHERE config_keys.id = $1;
