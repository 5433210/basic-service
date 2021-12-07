package apiv1

type Role struct {
	Id                  string            `json:"id"`
	Name                string            `json:"name"`
	Permissions         PermissionsInRole `json:"permissions"`
	Scopes              ScopesInRole      `json:"scopes"`
	ExclusiveRoles      []string          `json:"exclusive_roles"`
	DynamicallyIsolated bool              `json:"dynamically_isolated"`
}

type PermissionsInRole struct {
	Included map[string]interface{} `json:"included"`
	Excluded []string               `json:"excluded"`
}

type ScopesInRole struct {
	CanBeGranted  []string `json:"can_be_granted"`
	CanBeAccessed []string `json:"can_be_accessed"`
}

type RoleOptions struct {
	Id      string      `json:"id"`
	Options interface{} `json:"options"`
}
