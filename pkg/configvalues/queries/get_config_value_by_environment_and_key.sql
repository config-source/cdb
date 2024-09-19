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
INNER JOIN environments AS e ON cv.environment_id = e.id
INNER JOIN config_keys AS ck ON cv.config_key_id = ck.id
WHERE e.id = $1 AND ck.name = $2;
