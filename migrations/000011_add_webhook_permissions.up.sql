INSERT INTO permissions (name) VALUES
    ('CAN_MANAGE_WEBHOOKS');

-- Add permissions to the roles
INSERT INTO permissions_to_roles (permission_id, role_id)
SELECT permissions.id, roles.id
FROM roles
JOIN permissions
ON roles.name = 'Administrator' AND permissions.name = 'CAN_MANAGE_WEBHOOKS';