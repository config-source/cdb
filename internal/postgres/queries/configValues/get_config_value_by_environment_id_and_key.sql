SELECT
    cv.id,
    cv.environment_id,
    cv.config_key_id,
    ck.name,
    ck.value_type,
    cv.str_value,
    cv.int_value,
    cv.float_value,
    cv.bool_value,
    cv.created_at
FROM config_values AS cv 
INNER JOIN environments AS e ON config_values.environment_id = environments.id
INNER JOIN config_keys AS ck ON config_values.config_key_id = config_keys.id
WHERE cv.environment_id = $1 AND ck.name = $2;
