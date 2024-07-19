SELECT
    wh.*
FROM webhook_definitions AS wh
INNER JOIN webhook_definitions_to_environment AS wh_to_e ON wh.id = wh_to_e.webhook_defnintion_id
WHERE wh_to_e.environment_id = $1;