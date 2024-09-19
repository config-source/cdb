INSERT INTO environments (
    name,
    promotes_to_id,
    sensitive,
    service_id
) 
VALUES (
    $1, 
    $2,
    $3,
    $4
)
RETURNING *;
