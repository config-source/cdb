INSERT INTO webhook_definitions (
    template,
    url,
    authz_header
) 
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;