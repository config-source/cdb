SELECT config_keys.*, services.id as service_id, services.name as service_name FROM config_keys
INNER JOIN services ON services.id = config_keys.service_id;
