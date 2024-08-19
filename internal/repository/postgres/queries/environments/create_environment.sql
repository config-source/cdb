INSERT INTO environments (
    name,
    promotes_to_id
) 
VALUES (
    $1, 
    $2
)
RETURNING *;
