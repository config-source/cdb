SELECT environments.*, services.name as service_name FROM environments 
JOIN services ON services.id = environments.service_id
WHERE services.name = $1 AND environments.name = $2;
