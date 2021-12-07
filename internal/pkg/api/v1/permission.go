package apiv1

type Permission struct {
	Id      string      `json:"id"`
	Options interface{} `json:"options"`
}

type PermissionOptions struct {
	Id      string      `json:"id"`
	Options interface{} `json:"options"`
}
