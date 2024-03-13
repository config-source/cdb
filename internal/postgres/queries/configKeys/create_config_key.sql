INSERT INTO config_keys (
    name,
    value_type,
    can_propagate
) 
VALUES (
    $1, 
    $2,
    $3
)
RETURNING *;
