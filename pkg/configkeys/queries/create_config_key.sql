INSERT INTO config_keys (
    name,
    value_type,
    can_propagate,
    service_id
) 
VALUES (
    $1, 
    $2,
    $3,
    $4
)
RETURNING *;
