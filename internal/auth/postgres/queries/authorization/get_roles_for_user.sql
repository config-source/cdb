SELECT roles.name
FROM users
INNER JOIN users_to_roles ON (users.id = users_to_roles.user_id)
INNER JOIN roles ON (users_to_roles.role_id = roles.id)
WHERE users.id = $1;
