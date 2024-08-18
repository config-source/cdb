UPDATE config_values
SET environment_id = $1, 
    config_key_id  = $2,
    str_value      = $3,
    int_value      = $4,
    float_value    = $5,
    bool_value     = $6
WHERE id = $7
RETURNING *;
