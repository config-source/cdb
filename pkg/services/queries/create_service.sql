INSERT INTO services (
    name
) 
VALUES (
    $1
)
RETURNING *;
