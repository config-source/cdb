SELECT permissions.name
FROM permissions
INNER JOIN permissions_to_roles ON (permissions.id = permissions_to_roles.permission_id)
INNER JOIN roles ON (permissions_to_roles.role_id = roles.id)
WHERE roles.name = $1;
