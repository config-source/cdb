SELECT COUNT(*)
FROM users
INNER JOIN users_to_roles ON (users.id = users_to_roles.user_id)
INNER JOIN roles ON (users_to_roles.role_id = roles.id)
INNER JOIN permissions_to_roles ON (permissions_to_roles.role_id = roles.id)
INNER JOIN permissions ON (permissions.id = permissions_to_roles.permission_id)
WHERE users.id = $1 AND permissions.name = $2
LIMIT 1;
