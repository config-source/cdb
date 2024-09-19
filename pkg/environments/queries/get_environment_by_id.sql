SELECT environments.*, services.name as service_name FROM environments 
JOIN services ON services.id = environments.service_id
WHERE environments.id = $1;
