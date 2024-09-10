INSERT INTO environments (
    name,
    promotes_to_id,
    sensitive
) 
VALUES (
    $1, 
    $2,
    $3
)
RETURNING *;
