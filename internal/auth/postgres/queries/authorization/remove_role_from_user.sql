DELETE FROM users_to_roles
WHERE 
    user_id = $1 AND
    role_id IN (
        SELECT id
        FROM roles
        WHERE name = $2
    );
