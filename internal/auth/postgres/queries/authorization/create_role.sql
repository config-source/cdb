INSERT INTO roles (name)
VALUES ($1)
RETURNING *;