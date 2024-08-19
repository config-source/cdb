DELETE FROM roles_to_permissions
WHERE 
    role_id IN (
        SELECT id
        FROM roles
        WHERE roles.name = $1
    ) 
    AND permission_id IN (
        SELECT id
        FROM permissions
        WHERE permissions.name = $2
    );
