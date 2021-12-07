package apiv1

type Deny struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Permissions PermissionsInDeny `json:"permissions"`
	Scopes      ScopesInDeny      `json:"scopes"`
}

type PermissionsInDeny struct {
	Included []string `json:"included"`
	Excluded []string `json:"excluded"`
}

type ScopesInDeny struct {
	CanBeGranted  []string `json:"can_be_granted"`
	CanBeAccessed []string `json:"can_be_accessed"`
}

type DenyOptions struct {
	Id      string      `json:"id"`
	Options interface{} `json:"options"`
}
