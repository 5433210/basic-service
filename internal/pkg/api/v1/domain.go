package apiv1

type Domain struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Options interface{} `json:"options"`
}
