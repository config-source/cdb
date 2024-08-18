package auth

type Permission string

const (
	PermissionConfigureEnvironments          Permission = "CAN_CONFIGURE_ENVIRONMENTS"
	PermissionConfigureSensitiveEnvironments Permission = "CAN_CONFIGURE_SENSITIVE_ENVIRONMENTS"
	PermissionManageEnvironments             Permission = "CAN_MANAGE_ENVIRONMENTS"
	PermissionManageRoles                    Permission = "CAN_MANAGE_ROLES"
	PermissionManageUsers                    Permission = "CAN_MANAGE_USERS"
	PermissionManageConfigKeys               Permission = "CAN_MANAGE_CONFIG_KEYS"
)
