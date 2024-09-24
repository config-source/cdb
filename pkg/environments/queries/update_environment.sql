UPDATE environments
SET name = $2,
    promotes_to_id = $3,
    sensitive = $4
WHERE id = $1
RETURNING *;
