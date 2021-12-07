package apiv1

type Subject struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Options interface{} `json:"options"`
}

type SubjectOptions struct {
	Id      string      `json:"id"`
	Options interface{} `json:"options"`
}
