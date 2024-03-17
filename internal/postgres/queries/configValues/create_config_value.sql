INSERT INTO config_values (
    environment_id,
    config_key_id,
    str_value,
    int_value,
    float_value,
    bool_value
) 
VALUES (
    $1, 
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
