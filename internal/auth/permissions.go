package auth

type Permission string

const (
	PermissionConfigureEnvironments          Permission = "CAN_CONFIGURE_ENVIRONMENTS"
	PermissionConfigureSensitiveEnvironments            = "CAN_CONFIGURE_SENSITIVE_ENVIRONMENTS"
	PermissionManageEnvironments                        = "CAN_MANAGE_ENVIRONMENTS"
	PermissionManageRoles                               = "CAN_MANAGE_ROLES"
	PermissionManageUsers                               = "CAN_MANAGE_USERS"
	PermissionManageConfigKeys                          = "CAN_MANAGE_CONFIG_KEYS"
)
