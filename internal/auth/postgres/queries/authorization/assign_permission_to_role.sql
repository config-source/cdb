INSERT INTO roles_to_permissions (permission_id, role_id)
SELECT permissions.id, roles.id
FROM roles
-- TODO: gotta be a way to do multiple permissions at once since we get a list
-- but for now we just call this in a loop.
JOIN permissions ON (roles.name = $1 AND permissions.name = $2)
