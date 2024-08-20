DELETE FROM permissions_to_roles
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
