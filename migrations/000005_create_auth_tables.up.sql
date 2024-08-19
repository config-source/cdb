CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE CONSTRAINT email_not_empty CHECK (email <> ''),
    password TEXT NOT NULL UNIQUE CONSTRAINT password_not_empty CHECK (password <> '')
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CONSTRAINT name_not_empty CHECK (name <> '')
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CONSTRAINT name_not_empty CHECK (name <> '')
);

CREATE TABLE permissions_to_roles (
    permission_id integer REFERENCES permissions,
    role_id integer REFERENCES roles,
    UNIQUE (permission_id, role_id)
);

CREATE TABLE users_to_roles (
    user_id integer REFERENCES users,
    role_id integer REFERENCES roles,
    UNIQUE (user_id, role_id)
);

-- Insert the permissions
INSERT INTO permissions (name) VALUES
    ('CAN_CONFIGURE_ENVIRONMENTS'),
    ('CAN_CONFIGURE_SENSITIVE_ENVIRONMENTS'),
    ('CAN_MANAGE_ENVIRONMENTS'),
    ('CAN_MANAGE_USERS'),
    ('CAN_MANAGE_CONFIG_KEYS');

-- Insert the default roles
INSERT INTO roles (name) VALUES 
    ('Administrator'),
    ('Operator');

-- Add permissions to the roles
INSERT INTO permissions_to_roles (permission_id, role_id)
SELECT permissions.id, roles.id
FROM roles
JOIN permissions
-- Administrator gets all permissions so don't filter them.
ON roles.name = 'Administrator';

INSERT INTO permissions_to_roles (permission_id, role_id)
SELECT permissions.id, roles.id
FROM roles
JOIN permissions
ON (
    roles.name = 'Operator' AND
    (
        permissions.name = 'CAN_CONFIGURE_ENVIRONMENTS' OR
        permissions.name = 'CAN_CONFIGURE_SENSITIVE_ENVIRONMENTS'
    )
);
